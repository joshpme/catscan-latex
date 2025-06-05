# CatScan Latex

This is the LaTeX checker developed for [IPAC'25](https://ipac25.org/). The initial release perform a series of simple checks for common issues in bibitems. This is due to this section typically requiring the most attention by editors.

## Structure

1. `finder` is the document parser.
2. `checker` performs the detection of issues
3. `main` handles generating an output. Including generating a suitable summary to be used as the comment in indico. This uses google's Gemini AI agent.
4. `stats` is directly executable, for analysing the impact of changes against real world papers.

## Generating the baseline stats

Statistics are generated on example files to gauge the impact of changes to the checks.

To generate a new baseline (`stats/details.csv` and `stats/summary.csv`), run the following from the project root

```bash
go run stats/stats.go
```

You'll need some tex files in `examples/`

The recommendation is to run this before then after code changes to see the impact your changes make in the real world.