"""
Databricks OpenTelemetry Azure Monitoring

This package provides tools for instrumenting Databricks notebooks with OpenTelemetry
and integrating with Azure Application Insights for monitoring and observability.
"""

from .otel_helper import OpenTelemetryHelper

__all__ = ['OpenTelemetryHelper']
