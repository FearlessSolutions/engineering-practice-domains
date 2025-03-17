# Parent-Child Notebook Workflow with OpenTelemetry

```mermaid
flowchart TB
    subgraph Notebook_Workflow ["Notebook_Workflow (Parent Span)"]
        direction TB
        
        WorkflowStart["Start Notebook Workflow
            Generate Workflow ID"]
        
        subgraph Child1 ["Child_Notebook_1 Span"]
            direction TB
            Child1Start["Start Child Notebook 1 Execution"]
            Child1Process["Data Validation Process"]
            Child1Return["Return JSON Result
                {
                  status: 'Success',
                  total_records: 8500,
                  validation_errors: 25,
                  processing_time_sec: 3.45
                }"]
            Child1End["End Child Notebook 1 Span"]
            
            Child1Start --> Child1Process
            Child1Process --> Child1Return
            Child1Return --> Child1End
        end
        
        subgraph Child2 ["Child_Notebook_2 Span"]
            direction TB
            Child2Start["Start Child Notebook 2 Execution"]
            Child2Process["Data Aggregation Process"]
            Child2Return["Return JSON Result
                {
                  status_code: 200,
                  num_aggregations: 35,
                  execution_time_sec: 5.67,
                  memory_usage_mb: 350.5
                }"]
            Child2End["End Child Notebook 2 Span"]
            
            Child2Start --> Child2Process
            Child2Process --> Child2Return
            Child2Return --> Child2End
        end
        
        WorkflowMetrics["Calculate Workflow Metrics
            Total Records: 8500
            Total Errors: 25
            Total Duration: 10.23s"]
        WorkflowEnd["End Notebook Workflow
            Add Completion Event"]
        
        WorkflowStart --> Child1
        Child1 --> Child2
        Child2 --> WorkflowMetrics
        WorkflowMetrics --> WorkflowEnd
    end
    
    subgraph Implementation ["Implementation Methods"]
        direction TB
        
        RunNotebook["run_notebook_with_tracing Method
            - Automatic span creation
            - Execute notebook with dbutils.notebook.run
            - Parse JSON result
            - Set span attributes
            - Record metrics
            - Handle errors"]
        
        InstrumentFunction["instrument_function Method
            - Wrap any function with tracing
            - More flexible approach
            - Custom attribute mapping"]
    end
    
    subgraph Azure ["Azure Application Insights"]
        Traces["Distributed Traces"]
        Metrics["Custom Metrics"]
        
        Traces --> Dashboards["Dashboards & Visualizations"]
        Metrics --> Dashboards
    end
    
    Notebook_Workflow -- "Uses" --> RunNotebook
    Notebook_Workflow -- "Alternative" --> InstrumentFunction
    Notebook_Workflow -- "Export Telemetry" --> Azure
    
    classDef workflow fill:#1E88E5,color:white;
    classDef child fill:#43A047,color:white;
    classDef methods fill:#9C27B0,color:white;
    classDef azure fill:#0072C6,color:white;
    classDef metrics fill:#FFC107,color:black;
    
    class Notebook_Workflow workflow;
    class Child1,Child2 child;
    class Implementation,RunNotebook,InstrumentFunction methods;
    class Azure,Traces,Metrics,Dashboards azure;
    class WorkflowMetrics metrics;
```

This diagram illustrates the parent-child notebook workflow with OpenTelemetry instrumentation. It shows:

1. The hierarchical structure of spans:
   - Notebook_Workflow (Parent Span)
     - Child_Notebook_1 (Child Span)
     - Child_Notebook_2 (Child Span)

2. The execution flow of the notebook workflow:
   - Parent notebook starts and generates a workflow ID
   - Child Notebook 1 executes (data validation)
   - Child Notebook 2 executes (data aggregation)
   - Parent notebook calculates overall metrics
   - Parent notebook ends and adds completion event

3. The JSON results returned by child notebooks:
   - Child Notebook 1 returns validation results
   - Child Notebook 2 returns aggregation results

4. The implementation methods provided by OpenTelemetryHelper:
   - `run_notebook_with_tracing`: Simplified method for notebook execution with automatic tracing
   - `instrument_function`: More flexible method for wrapping any function with tracing

5. The export of telemetry data to Azure Application Insights:
   - Traces for distributed tracing
   - Custom metrics for performance monitoring
   - Dashboards for visualization

The color coding helps distinguish between:
- Blue: Overall notebook workflow
- Green: Child notebook executions
- Purple: Implementation methods
- Yellow: Metrics and performance data
- Azure Blue: Azure monitoring components
