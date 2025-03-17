# ETL Pipeline with OpenTelemetry Visualization

```mermaid
flowchart LR
    subgraph ETL_Pipeline ["ETL_Pipeline (Parent Span)"]
        direction TB
        
        subgraph Extraction ["DataExtraction Span"]
            direction TB
            Extract["Extract Data from API"]
            ExtractMetrics["Records Extracted: 15000
                HTTP Status: 200
                Duration: 2.5s"]
            
            Extract --> ExtractMetrics
        end
        
        subgraph Transformation ["DataTransformation Span"]
            direction TB
            Transform["Transform Data"]
            TransformMetrics["Records Input: 15000
                Records Transformed: 15000
                Records Failed: 0
                Duration: 1.8s
                Rate: 8333 records/sec"]
            
            Transform --> TransformMetrics
        end
        
        subgraph Loading ["DataLoading Span"]
            direction TB
            Load["Load Data to Storage"]
            LoadMetrics["Records Written: 15000
                Storage Path: /mnt/output/transformed_data
                Duration: 3.2s
                Rate: 4687 records/sec"]
            
            Load --> LoadMetrics
        end
        
        PipelineStart["Start ETL Pipeline
            Generate Pipeline ID"]
        PipelineMetrics["ETL Efficiency: 100%
            Total Duration: 7.5s
            Total Records: 15000"]
        PipelineEnd["End ETL Pipeline
            Add Completion Event"]
        
        PipelineStart --> Extraction
        Extraction --> Transformation
        Transformation --> Loading
        Loading --> PipelineMetrics
        PipelineMetrics --> PipelineEnd
    end
    
    subgraph Azure ["Azure Application Insights"]
        Traces["Distributed Traces"]
        Metrics["Custom Metrics"]
        Logs["Application Logs"]
        
        Traces --> Dashboards["Dashboards & Visualizations"]
        Metrics --> Dashboards
        Logs --> Dashboards
        
        Dashboards --> Alerts["Alerts & Notifications"]
    end
    
    ETL_Pipeline -- "Export Telemetry" --> Azure
    
    classDef pipeline fill:#1E88E5,color:white;
    classDef stage fill:#43A047,color:white;
    classDef azure fill:#0072C6,color:white;
    classDef metrics fill:#FFC107,color:black;
    
    class ETL_Pipeline pipeline;
    class Extraction,Transformation,Loading stage;
    class Azure,Traces,Metrics,Logs,Dashboards,Alerts azure;
    class ExtractMetrics,TransformMetrics,LoadMetrics,PipelineMetrics metrics;
```

This diagram visualizes the ETL pipeline with OpenTelemetry instrumentation. It shows:

1. The hierarchical structure of spans:
   - ETL_Pipeline (Parent Span)
     - DataExtraction (Child Span)
     - DataTransformation (Child Span)
     - DataLoading (Child Span)

2. The flow of data through the pipeline stages:
   - Extraction: Retrieving data from an API
   - Transformation: Processing and cleaning the data
   - Loading: Writing the transformed data to storage

3. Key metrics captured at each stage:
   - Record counts
   - Processing durations
   - Processing rates
   - Status codes and error information

4. The export of telemetry data to Azure Application Insights:
   - Traces for distributed tracing
   - Custom metrics for performance monitoring
   - Logs for error tracking
   - Dashboards and alerts for visualization and notification

The color coding helps distinguish between:
- Blue: Overall ETL pipeline
- Green: Individual pipeline stages
- Yellow: Metrics and performance data
- Azure Blue: Azure monitoring components
