import os
import time
import json
import functools
from azure.monitor.opentelemetry.exporter import AzureMonitorTraceExporter, AzureMonitorMetricExporter
from databricks.sdk.runtime import dbutils
from opentelemetry import trace, metrics
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.sdk.metrics import MeterProvider
from opentelemetry.sdk.metrics.export import PeriodicExportingMetricReader

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
            
    def _trace_execution(self, span_name, func, func_args=None, func_kwargs=None, 
                        attributes_mapping=None, additional_attributes=None, 
                        pre_execution_callback=None, post_execution_callback=None):
        """
        Private helper method to handle common tracing logic.
        
        :param span_name: Name of the span to create
        :param func: Function to execute
        :param func_args: Arguments to pass to the function (tuple)
        :param func_kwargs: Keyword arguments to pass to the function (dict)
        :param attributes_mapping: Dictionary mapping span attribute names to function return dict keys
        :param additional_attributes: Additional attributes to set on the span
        :param pre_execution_callback: Function to call before executing func (receives span_name)
        :param post_execution_callback: Function to call after executing func but before ending span (receives span_name, result)
        :return: Result of the function execution
        """
        func_args = func_args or ()
        func_kwargs = func_kwargs or {}
        additional_attributes = additional_attributes or {}
        
        # Start tracing with base attributes
        span_attributes = {"etl_pipeline_id": self.etl_pipeline_id}
        span_attributes.update(additional_attributes)
        
        self.start_tracing(span_name, span_attributes)
        
        # Set function name as a span attribute if available
        if hasattr(func, "__name__"):
            self.set_span_attribute(span_name, "function_name", func.__name__)
            print(f"Started tracing for function '{func.__name__}' with span '{span_name}'")
        
        # Execute pre-execution callback if provided
        if pre_execution_callback:
            pre_execution_callback(span_name)
        
        result = None
        error = None
        
        try:
            # Execute the function
            result = func(*func_args, **func_kwargs)
            
            # Process the result
            if isinstance(result, dict):
                # Set span attributes from the function's return value
                if attributes_mapping:
                    for attr_name, result_key in attributes_mapping.items():
                        if result_key in result:
                            self.set_span_attribute(span_name, attr_name, result[result_key])
                
                # Record metrics if applicable
                if span_name in self.metrics_registry:
                    for metric_name in self.metrics_registry[span_name]:
                        if metric_name in result:
                            # Add type checking before recording the metric
                            if isinstance(result[metric_name], (int, float)):
                                self.record_metric(span_name, metric_name, result[metric_name])
                            else:
                                print(f"Warning: Metric '{metric_name}' value is not a number, skipping")
            
            return result
            
        except Exception as e:
            # Set error attributes
            error_msg = str(e)
            self.set_span_attribute(span_name, "error", "true")
            self.set_span_attribute(span_name, "error_message", error_msg)
            error = e
            raise
            
        finally:
            # Execute post-execution callback if provided (before ending the span)
            processed_result = None
            if post_execution_callback and result is not None:
                try:
                    processed_result = post_execution_callback(span_name, result)
                except Exception as callback_error:
                    # Log the callback error but don't override the original error if there was one
                    print(f"Error in post-execution callback: {str(callback_error)}")
                    if error is None:
                        # Only set error attributes if there wasn't already an error
                        self.set_span_attribute(span_name, "error", "true")
                        self.set_span_attribute(span_name, "error_message", f"Post-execution callback error: {str(callback_error)}")
            
            # End tracing
            self.end_tracing(span_name)
            if hasattr(func, "__name__"):
                print(f"Ended tracing for function '{func.__name__}' with span '{span_name}'")
    
    def instrument_function(self, function, span_name=None, attributes_mapping=None):
        """
        Wraps a function with OpenTelemetry tracing.
        
        :param function: The function to wrap
        :param span_name: Name for the span (defaults to function name)
        :param attributes_mapping: Dictionary mapping span attribute names to function return dict keys
        :return: Wrapped function
        
        Example usage:
        
        run_traced_notebook = workflow_otel_helper.instrument_function(
            dbutils.notebook.run,
            span_name="Child_Notebook_2"
        )
        
        child2_result_json = run_traced_notebook("./child_notebook_2", timeout_seconds=600)
        """
        # Use function name as default span name if not provided
        if span_name is None:
            span_name = function.__name__
        
        outer_self = self
        
        @functools.wraps(function)
        def trace_wrapper(*args, **kwargs):
            return outer_self._trace_execution(
                span_name=span_name,
                func=function,
                func_args=args,
                func_kwargs=kwargs,
                attributes_mapping=attributes_mapping
            )
        
        return trace_wrapper
    
    def run_notebook_with_tracing(self, notebook_path, span_name=None, timeout_seconds=600, arguments=None, **kwargs):
        """
        Wrapper around dbutils.notebook.run that automatically handles tracing.
        
        :param notebook_path: Path to the notebook to run
        :param span_name: Name for the span (defaults to notebook name if not provided)
        :param timeout_seconds: Timeout for notebook execution
        :param arguments: Arguments to pass to the notebook
        :param kwargs: Additional keyword arguments to include as span attributes
        :return: The parsed result from the notebook execution (as a dictionary if JSON, otherwise as string)
        
        Example usage:
        
        child1_result = workflow_otel_helper.run_notebook_with_tracing(
            "./child_notebook_1", 
            span_name="Child_Notebook_1",
            timeout_seconds=600,
            notebook_type="validation"
        )
        """
        # Extract notebook name from path if span_name not provided
        if span_name is None:
            span_name = notebook_path.split('/')[-1].replace('.', '_')
        
        # Define notebook-specific callbacks
        def pre_execution(span_name):
            self.add_span_event(span_name, "Notebook Execution Started", {
                "notebook_path": notebook_path,
                "timestamp": time.time()
            })
        
        def post_execution(span_name, result_json):
            # Add events but don't try to parse JSON here
            # We'll parse it after the span is ended
            self.add_span_event(span_name, "Notebook Execution Completed", {
                "timestamp": time.time()
            })
        
        # Define error handler for notebook execution
        def notebook_error_handler(e):
            error_msg = str(e)
            self.add_span_event(span_name, "Notebook Execution Failed", {
                "error_message": error_msg,
                "timestamp": time.time()
            })
        
        # Execute notebook with tracing
        try:
            result_json = self._trace_execution(
                span_name=span_name,
                func=dbutils.notebook.run,
                func_args=(notebook_path, timeout_seconds, arguments),
                additional_attributes={"notebook_path": notebook_path, **kwargs},
                pre_execution_callback=pre_execution,
                post_execution_callback=post_execution
            )
            
            # Now parse the JSON after the span has been ended
            try:
                # Parse the result if it's JSON
                result = json.loads(result_json)
                
                # Return the parsed result
                return result
            except json.JSONDecodeError:
                # If result is not JSON, just return it as is
                return result_json
                
        except Exception as e:
            # This will only be called if an error occurs that wasn't caught by _trace_execution
            notebook_error_handler(e)
            raise
    
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
                return self._trace_execution(
                    span_name=span_name,
                    func=func,
                    func_args=args,
                    func_kwargs=kwargs,
                    attributes_mapping=attributes_mapping
                )
            return wrapper
        return decorator
