# MySQL InnoDB File Parser

A tool to parse MySQL InnoDB files directly.

## Quick Start


### Setup

1. Start the environment:

```powershell
docker-compose up -d
```


2. Watch the progress:

```powershell
docker-compose logs -f mysql
```

You'll see messages like:

```
Created 100000 numbers...
Created 200000 numbers...
...
Shuffling and inserting 1,000,000 rows...
```

3. Check when ready:

```powershell
# Quick check
docker-compose exec mysql mysql -uroot test -e "SELECT COUNT(*) FROM t;"

```

4. Once data is loaded, test the app:

```powershell
docker-compose exec app bash
go build -o /bin/goParseInnodb ./cmd/goParseInnodb
/bin/goParseInnodb
```

### Verify Setup

```powershell
# Check row count (should be 1,000,000)
docker-compose exec mysql mysql -uroot test -e "SELECT COUNT(*) FROM t;"

# See sample data
docker-compose exec mysql mysql -uroot test -e "SELECT * FROM t LIMIT 10;"

# Check file size
docker-compose exec app ls -lh /mysql-data/test/t.ibd
```

### Development Workflow

```powershell
# Make changes to cmd/goParseInnodb/main.go
# Then rebuild and run:
docker-compose exec app bash
go build -o /bin/goParseInnodb ./cmd/goParseInnodb
/bin/goParseInnodb

# Or one-liner from PowerShell:
docker-compose exec app bash -c "go build -o /bin/goParseInnodb ./cmd/goParseInnodb && /bin/goParseInnodb"
```

## Project Structure

- `cmd/goParseInnodb/` - Main application
- `pkg/` - Reusable packages (future)
- `docker/` - Docker configuration
- `output/` - Output files

## References
- `https://blog.jcole.us/innodb/` - Amazing set of articles about innodb internal architecture and data representation







