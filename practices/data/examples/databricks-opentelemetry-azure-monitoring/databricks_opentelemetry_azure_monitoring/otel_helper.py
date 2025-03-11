import os
import time
import json
import uuid
import functools
from opentelemetry import trace, metrics
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.sdk.metrics import MeterProvider
from opentelemetry.sdk.metrics.export import PeriodicExportingMetricReader
from azure.monitor.opentelemetry.exporter import AzureMonitorTraceExporter, AzureMonitorMetricExporter

class OpenTelemetryHelper:
    """Encapsulates OpenTelemetry tracing, trace attributes, and metrics integration.
    
    This helper class provides methods for:
    - Starting and ending trace spans
    - Setting span attributes
    - Recording metrics
    - Adding span events
    - Decorating functions with automatic tracing
    """

    def __init__(self, span_name, etl_pipeline_id, span_metrics=None, span_attributes=None):
        """
        Initializes OpenTelemetry tracing with a parent span, metrics, and attributes.

        :param span_name: Name of the parent span for tracing.
        :param etl_pipeline_id: Unique ID for the ETL pipeline.
        :param span_metrics: (Optional) Dictionary where keys are span names and values are metric definitions.
        :param span_attributes: (Optional) Dictionary of initial trace attributes for the parent span.
        """

        # Store ETL metadata
        self.etl_pipeline_id = etl_pipeline_id
        self.span_metrics = span_metrics or {}  
        self.span_attributes = span_attributes or {}  

        # Setup OpenTelemetry tracing
        trace.set_tracer_provider(TracerProvider())
        self.tracer_provider = trace.get_tracer_provider()
        self.tracer = trace.get_tracer(__name__)
        azure_exporter = AzureMonitorTraceExporter.from_connection_string(os.environ['AZURE_MONITOR_CONNECTION_STRING'])
        self.tracer_provider.add_span_processor(BatchSpanProcessor(azure_exporter))

        # Setup OpenTelemetry metrics
        metric_exporter = AzureMonitorMetricExporter.from_connection_string(os.environ['AZURE_MONITOR_CONNECTION_STRING'])
        metric_reader = PeriodicExportingMetricReader(metric_exporter)
        metrics.set_meter_provider(MeterProvider(metric_readers=[metric_reader]))
        self.meter = metrics.get_meter(__name__)

        # Initialize metric storage
        self.metrics_registry = {}
        
        # Initialize tracing storage
        self.active_spans = {}
        
        # Track span start times for duration calculation
        self.span_start_times = {}

        # Start the parent span
        self.start_tracing(span_name)

        # Assign predefined attributes to parent span
        for key, value in self.span_attributes.items():
            self.active_spans[span_name].set_attribute(key, value)

        # Initialize metrics for spans
        self.config_metrics(self.span_metrics)

    def start_tracing(self, span_name, attributes=None):
        """
        Starts a tracing span, useful for wrapping around notebook execution.

        :param span_name: Name of the tracing span.
        :param attributes: (Optional) Dictionary of attributes to set on the span.
        """
        if span_name in self.active_spans:
            raise ValueError(f"Span '{span_name}' is already active!")

        span = self.tracer.start_span(span_name)
        self.active_spans[span_name] = span
        
        # Record the start time for duration calculation
        self.span_start_times[span_name] = time.time()
        
        print(f"Started Tracing Span: {span_name}")

        if attributes:
            for key, value in attributes.items():
                span.set_attribute(key, value)

    def end_tracing(self, span_name):
        """
        Ends a tracing span.

        :param span_name: Name of the tracing span to end.
        """
        if span_name in self.active_spans:
            # Calculate duration if we have a start time
            if span_name in self.span_start_times:
                duration_sec = time.time() - self.span_start_times[span_name]
                # Add duration as an attribute (in seconds, rounded to milliseconds precision)
                self.active_spans[span_name].set_attribute("duration_sec", round(duration_sec, 3))
                print(f"Span duration for '{span_name}': {round(duration_sec, 3)} seconds")
                # Clean up the start time
                del self.span_start_times[span_name]
            
            self.active_spans[span_name].end()
            print(f"Ended Tracing Span: {span_name}")
            del self.active_spans[span_name]
        else:
            raise KeyError(f"Tracing span '{span_name}' not found!")

    def config_metrics(self, span_metrics):
        """
        Configures metrics for spans.

        :param span_metrics: Dictionary where keys are span names and values are metric definitions.
        """
        for span_name, metrics_def in span_metrics.items():
            if span_name not in self.metrics_registry:
                self.metrics_registry[span_name] = {}

            for metric_name, metric_type in metrics_def.items():
                if metric_type == "counter":
                    self.metrics_registry[span_name][metric_name] = {
                        "metric": self.meter.create_counter(metric_name),
                        "type": "counter"
                    }
                elif metric_type == "histogram":
                    self.metrics_registry[span_name][metric_name] = {
                        "metric": self.meter.create_histogram(metric_name),
                        "type": "histogram"
                    }
                else:
                    raise ValueError(f"Unsupported metric type '{metric_type}' for metric '{metric_name}'.")

    def record_metric(self, span_name, metric_name, value, attributes=None):
        """
        Records a metric value for a specific span.

        :param span_name: Name of the span to associate the metric with.
        :param metric_name: Name of the metric to record.
        :param value: Value to add to the metric.
        :param attributes: Optional attributes for the metric.
        """
        if span_name in self.metrics_registry and metric_name in self.metrics_registry[span_name]:
            metric_info = self.metrics_registry[span_name][metric_name]
            metric_attrs = attributes or {"etl_pipeline_id": self.etl_pipeline_id}
            
            if metric_info["type"] == "counter":
                metric_info["metric"].add(value, metric_attrs)
            elif metric_info["type"] == "histogram":
                metric_info["metric"].record(value, metric_attrs)
        else:
            raise KeyError(f"Metric '{metric_name}' for span '{span_name}' is not defined.")

    def set_span_attribute(self, span_name, attribute_name, attribute_value):
        """
        Sets an attribute on a specific span.

        :param span_name: Name of the span to set the attribute on.
        :param attribute_name: Name of the attribute to set.
        :param attribute_value: Value to set for the attribute.
        """
        if span_name in self.active_spans:
            self.active_spans[span_name].set_attribute(attribute_name, attribute_value)
        else:
            raise KeyError(f"Tracing span '{span_name}' not found!")

    def add_span_event(self, span_name, event_name, attributes=None):
        """
        Adds an event to a specific span.

        :param span_name: Name of the span to add the event to.
        :param event_name: Name of the event to add.
        :param attributes: (Optional) Dictionary of attributes for the event.
        """
        if span_name in self.active_spans:
            self.active_spans[span_name].add_event(event_name, attributes)
            print(f"Added event '{event_name}' to span '{span_name}'")
        else:
            raise KeyError(f"Tracing span '{span_name}' not found!")
            
    def trace_function(self, span_name, attributes_mapping=None):
        """
        Decorator for tracing function execution.
        
        :param span_name: Name of the span to create.
        :param attributes_mapping: Dictionary mapping span attribute names to function return dict keys.
        :return: Decorated function.
        
        Example usage:
        
        @etl_otel_helper.trace_function("DataLoading", {
            "records_written": "records_written",
            "storage_path": "storage_path"
        })
        def perform_data_loading():
            # Function implementation
            return {
                "records_written": 100,
                "storage_path": "/path/to/data"
            }
        """
        def decorator(func):
            @functools.wraps(func)
            def wrapper(*args, **kwargs):
                # Start the span
                self.start_tracing(span_name, {"etl_pipeline_id": self.etl_pipeline_id})
                print(f"Started tracing for function '{func.__name__}' with span '{span_name}'")
                
                try:
                    # Execute the function
                    result = func(*args, **kwargs)
                    
                    # Set span attributes from the function's return value
                    if attributes_mapping and isinstance(result, dict):
                        for attr_name, result_key in attributes_mapping.items():
                            if result_key in result:
                                self.set_span_attribute(span_name, attr_name, result[result_key])
                    
                    # Record metrics if applicable
                    if span_name in self.metrics_registry:
                        for metric_name in self.metrics_registry[span_name]:
                            if metric_name in result:
                                self.record_metric(span_name, metric_name, result[metric_name])
                    
                    return result
                except Exception as e:
                    # Set error attributes
                    self.set_span_attribute(span_name, "error", "true")
                    self.set_span_attribute(span_name, "error_message", str(e))
                    raise
                finally:
                    # End the span
                    self.end_tracing(span_name)
                    print(f"Ended tracing for function '{func.__name__}' with span '{span_name}'")
            
            return wrapper
        return decorator
