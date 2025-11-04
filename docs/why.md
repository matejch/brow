# Why Not MCP?

Simple CLI tools can be more efficient than Model Context Protocol (MCP) servers for browser automation.

## Key Advantages

### 1. Token Efficiency
- **brow documentation**: ~200 tokens (this README)
- **MCP server**: 13,000-18,000 tokens for tool descriptions
- **Savings**: 98%+ reduction in context overhead

### 2. Composability
```bash
# Results can be saved to files and chained
./brow eval 'getData()' > data.json
cat data.json | jq '.items[]' | while read item; do
  # Process each item
done
```

With MCP, results flow through agent context, limiting composition.

### 3. Simplicity
- **brow**: Standard CLI with `--help` flags
- **MCP**: Requires server lifecycle, protocol understanding, separate process

### 4. Extensibility
Adding a new tool to brow:
1. Create `cmd/newtool.go`
2. Implement the command
3. Done

Adding to MCP:
1. Modify server code
2. Update tool definitions
3. Restart server
4. Update documentation

### 5. Bash-Native
Agents already excel at Bash. Why add another layer?

```bash
# This is natural for agents
./brow nav https://example.com
./brow screenshot > image.png

# This requires protocol overhead
mcp-call browser navigate https://example.com
mcp-call browser screenshot --output image.png
```

## When to Use MCP

MCP excels when you need:
- Complex stateful interactions across many tools
- Fine-grained permissions and security boundaries
- Integration with non-CLI tools
- Standardized cross-platform tool discovery

## When to Use brow

brow excels when you need:
- Browser automation for AI agents
- Minimal token overhead
- File-based composition
- Simple extension and customization
- Standard Unix tool philosophy

## Token Cost Comparison (Real Example)

### Scraping Quotes from quotes.toscrape.com

**With MCP:**
- Initial context: 13,000 tokens (tool descriptions)
- Per-operation: ~500 tokens (request + response through context)
- Total for 10 operations: ~18,000 tokens

**With brow:**
- Initial context: 200 tokens (README)
- Per-operation: ~50 tokens (command text)
- Results saved to files: 0 tokens in context
- Total for 10 operations: ~700 tokens

**Savings: ~96%**

## Design Principles

1. **Each tool does one thing well** - Unix philosophy
2. **Text-based output** - Easy for agents to parse
3. **File-based state** - Results persist outside context
4. **Minimal documentation** - Tools are self-explanatory
5. **Stateful browser** - Chrome maintains state between calls
6. **Zero protocol overhead** - Just stdin/stdout

## The Bottom Line

> "Agents can run Bash and write code well. Bash and code are composable."

If your task fits the "run command, get output" model, simple CLI tools will beat protocol-based approaches in:
- Token efficiency
- Simplicity
- Composability
- Development speed

MCP is powerful, but not always necessary. Choose the right tool for the job.
