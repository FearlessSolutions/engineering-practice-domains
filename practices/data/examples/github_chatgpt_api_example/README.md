# Scraping Github with ChatGPT: How AI helped me select software for my customer

This code was the subject of the COI-AI talk I gave on Tuesday, Sept 4, 2024. This was a personal project of mine to familiarize myself with the OpenAI API and do a bit of market research into geospatial software in terms of what are the most important geospatial software, both open source and commercial.

NOTE: This technique can be used to gain information about software from ANY domain, not just geospatial. You can simply replace the "geospatial" keyword in the python code with "AI" "machine learning" "databases", etc to gain insights and latest trends.

Goal of this repo: Get up to speed with the latest and greatest open source projects in the geospatial industry space SO THAT WE CAN MAKE DECISIONS.

Tools Used: Python, Github Repository Search API, OpenAI Python library Assistants API (talks to ChatGPT) , MongoDB to store JSON data

## Workflow

1. Do a Github repository search for a certain keyword, e.g., “geospatial” (query_github_openai.ipynb, part 1)
2. For each repository, ask ChatGPT questions about it (query_github_openai.ipynb, part 2)
3. Visualize AI-produced results (openai_data_exploratory.ipynb)




