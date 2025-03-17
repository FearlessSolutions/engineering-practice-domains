# Glossary of Terms

This glossary provides definitions for technical terms used throughout the documentation to help users better understand the concepts and terminology related to OpenTelemetry, Azure Application Insights, and Databricks.

## Table of Contents

- [Azure Application Insights](#azure-application-insights)
- [Azure Monitor](#azure-monitor)
- [Batch Span Processor](#batch-span-processor)
- [Child Span](#child-span)
- [Counter](#counter)
- [Databricks](#databricks)
- [ETL Pipeline](#etl-pipeline)
- [Histogram](#histogram)
- [KQL (Kusto Query Language)](#kql-kusto-query-language)
- [Metric](#metric)
- [OpenTelemetry](#opentelemetry)
- [Parent Span](#parent-span)
- [Parent-Child Notebook](#parent-child-notebook)
- [Span](#span)
- [Span Attribute](#span-attribute)
- [Span Event](#span-event)
- [Telemetry](#telemetry)
- [Trace](#trace)
- [Tracer](#tracer)

## Azure Application Insights

A feature of Azure Monitor that provides application performance monitoring (APM) capabilities. It collects telemetry data from applications and allows you to analyze it for performance monitoring, diagnostics, and analytics.

## Azure Monitor

A comprehensive monitoring solution for collecting, analyzing, and acting on telemetry from cloud and on-premises environments. It helps you maximize the availability and performance of your applications and services.

## Batch Span Processor

A component in OpenTelemetry that processes spans in batches before exporting them, which is more efficient than processing them individually.

## Child Span

A span that is a descendant of a parent span. Child spans represent sub-operations within a larger operation. For example, in an ETL pipeline, the extraction, transformation, and loading stages would be represented as child spans of the parent ETL pipeline span.

## Counter

A type of metric that represents a single numerical value that only increases over time. Examples include the number of records processed, errors encountered, or API calls made.

## Databricks

A cloud-based data engineering platform that provides a collaborative environment for data science and engineering teams to work with big data and machine learning.

## ETL Pipeline

Extract, Transform, Load - a process that extracts data from a source, transforms it to fit operational needs, and loads it into a destination system. ETL pipelines are commonly used for data integration and data warehousing.

## Histogram

A type of metric that tracks the distribution of values, such as request durations or response sizes. Histograms allow you to analyze the distribution of values, including percentiles, mean, and standard deviation.

## KQL (Kusto Query Language)

The query language used in Azure Application Insights to retrieve and analyze telemetry data. KQL is designed for data exploration and provides powerful filtering, aggregation, and visualization capabilities.

## Metric

A measurement of a specific aspect of a system's behavior over time. Metrics are typically numeric values that can be aggregated, such as counts, rates, durations, or sizes.

## OpenTelemetry

An open-source observability framework that provides a collection of tools, APIs, and SDKs to instrument, generate, collect, and export telemetry data (metrics, logs, and traces) for analysis.

## Parent Span

A span that encompasses one or more child spans. Parent spans represent higher-level operations, such as an entire ETL pipeline or a notebook workflow.

## Parent-Child Notebook

A pattern in Databricks where a parent notebook orchestrates the execution of one or more child notebooks. The parent notebook can pass parameters to child notebooks and collect their results.

## Span

The basic unit of work in OpenTelemetry tracing. A span represents a single operation within a trace, such as a function call, a database query, or an HTTP request. Spans have a start time, an end time, and can contain attributes, events, and links to other spans.

## Span Attribute

Key-value pairs that provide additional context about a span. Attributes can include information such as the operation name, status code, error details, or custom metadata.

## Span Event

A timestamped annotation within a span that marks a significant occurrence during the span's lifetime. Events can include information such as exceptions, state changes, or checkpoints.

## Telemetry

Data collected from remote systems for monitoring and analysis. Telemetry includes metrics, logs, and traces that provide insights into the behavior, performance, and health of applications and infrastructure.

## Trace

A collection of spans that form a tree-like structure representing the path of a request through a distributed system. Traces provide end-to-end visibility into the execution of a request across multiple services or components.

## Tracer

A component in OpenTelemetry that creates and manages spans. Tracers are responsible for starting and ending spans, adding attributes and events, and establishing parent-child relationships between spans.
