# HIBP (Have I Been Pwned) Command

The `/hibp` command allows you to check if a phone number or identifier has been exposed in data breaches, with a focus on the HiTeckGroop.in database.

## Prerequisites

To use this command, you need to obtain an API token from the Have I Been Pwned service:

1. Message the bot with `/api` to get your personal token
2. Initially, you'll have 100 free requests for testing
3. After that, requests will be charged from your balance

## Configuration

Set the `HIBP_TOKEN` environment variable with your API token:

```bash
export HIBP_TOKEN="your_hibp_api_token_here"
```

Or in your `.env` file:
```
HIBP_TOKEN=your_hibp_api_token_here
```

Or in docker-compose.yml:
```yaml
environment:
  HIBP_TOKEN: ${HIBP_TOKEN}
```

## Usage

The command is owner-only for security reasons:

```
/hibp <phone_or_identifier>
```

### Examples

```
/hibp 917888313823
/hibp anuragbpre018@gmail.com
```

## Response Format

The command will return:

1. API usage information (requests left, price, search time)
2. Data specifically from HiTeckGroop.in database if found:
   - Full name
   - Father's name
   - Document number
   - Addresses
   - Region information
   - All associated phone numbers
3. Summary of data found in other databases
4. Security recommendations

## Security Considerations

- This command is restricted to bot owners only
- All API requests are made over HTTPS
- No sensitive data is logged by the bot
- Results are only visible to the command issuer

## Pricing

The price of each request depends on the type and search limit:

**Formula:** `(5 + sqrt(Limit * Complexity)) / 5000`

Where:
- **Limit** is the search limit (default 100)
- **Complexity** depends on the number of words in your request:
  - 1 word: Complexity = 1
  - 2 words: Complexity = 5
  - 3 words: Complexity = 16
  - More than 3 words: Complexity = 40

With the default limit of 100, most requests cost approximately $0.003.