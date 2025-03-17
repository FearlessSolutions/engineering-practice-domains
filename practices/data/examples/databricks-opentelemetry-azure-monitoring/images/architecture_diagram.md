# OpenTelemetry Architecture Diagram

```mermaid
flowchart TB
    subgraph Databricks ["Databricks Workspace"]
        subgraph Notebooks ["Databricks Notebooks"]
            ETL["ETL Pipeline Notebook"]
            Parent["Parent Notebook"]
            Child1["Child Notebook 1"]
            Child2["Child Notebook 2"]
            
            Parent --> Child1
            Parent --> Child2
        end
        
        subgraph OTel ["OpenTelemetry Instrumentation"]
            Helper["OpenTelemetryHelper Class"]
            Tracer["OpenTelemetry Tracer"]
            Meter["OpenTelemetry Meter"]
            Exporter["Azure Monitor Exporter"]
            
            Helper --> Tracer
            Helper --> Meter
            Tracer --> Exporter
            Meter --> Exporter
        end
        
        ETL --> Helper
        Parent --> Helper
    end
    
    subgraph Azure ["Azure Cloud"]
        AppInsights["Azure Application Insights"]
        Dashboards["Azure Dashboards"]
        Alerts["Azure Alerts"]
        
        AppInsights --> Dashboards
        AppInsights --> Alerts
    end
    
    Exporter --> AppInsights
    
    classDef databricks fill:#1E88E5,color:white;
    classDef otel fill:#43A047,color:white;
    classDef azure fill:#0072C6,color:white;
    
    class Databricks,Notebooks databricks;
    class OTel,Helper,Tracer,Meter,Exporter otel;
    class Azure,AppInsights,Dashboards,Alerts azure;
```

This diagram illustrates the architecture of the Databricks OpenTelemetry Azure Monitoring solution. It shows how:

1. Databricks notebooks (ETL Pipeline and Parent-Child workflows) use the OpenTelemetryHelper class
2. The OpenTelemetryHelper class manages the OpenTelemetry Tracer and Meter
3. Telemetry data is exported to Azure Application Insights via the Azure Monitor Exporter
4. Azure Application Insights provides dashboards and alerts for monitoring

The color coding helps distinguish between:
- Blue: Databricks components
- Green: OpenTelemetry components
- Azure Blue: Azure components
