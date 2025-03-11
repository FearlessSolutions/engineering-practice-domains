from setuptools import setup, find_packages

setup(
    name="databricks-opentelemetry-azure-monitoring",
    version="0.1.0",
    description="OpenTelemetry instrumentation for Databricks notebooks with Azure Monitor integration",
    author="Your Organization",
    author_email="your.email@example.com",
    packages=find_packages(),
    install_requires=[
        "azure-monitor-opentelemetry>=1.0.0",
        "opentelemetry-sdk>=1.12.0",
        "azure-core>=1.24.0",
        "opentelemetry-api>=1.12.0",
        "opentelemetry-instrumentation>=0.33b0",
    ],
    python_requires=">=3.6",
)
