# TODO

## Performance

### Query optimalisatie pagination / sync concurrency

Tijdens sync (goroutine) moet pagination wachten. Waarschijnlijk geen PostgreSQL lock
(MVCC laat reads en writes naast elkaar lopen), maar mogelijk:

- `PrepareStmt: true` in GORM-config veroorzaakt mutex-contention op de statement-cache
- Geen connection pool limiet ingesteld (`SetMaxOpenConns` / `SetMaxIdleConns`)
- `SaveEvents` draait honderden losse impliciete transacties sequentieel per event

**Te onderzoeken:**
- Draai `SELECT pid, state, wait_event_type, wait_event, query FROM pg_stat_activity WHERE state != 'idle'` tijdens sync om te zien of er echte PostgreSQL locks zijn
- Overweeg `SetMaxOpenConns(20)` / `SetMaxIdleConns(10)` toe te voegen aan `Storage.Connect()`
- Bekijk of `SaveEvents`-loop gebundeld kan worden in één transactie (bulk insert) i.p.v. per event
- ~~Index op `notes.event_created_at`~~ — aangemaakt via migratie 000006 ✓
