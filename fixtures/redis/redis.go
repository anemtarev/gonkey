package redis

import (
    "context"

    "github.com/go-redis/redis/v9"
    "github.com/lamoda/gonkey/fixtures/redis/parser"
)

type loader struct {
    locations []string
    client    *redis.Client
}

type LoaderOptions struct {
    FixtureDir string
    Redis      *redis.Options
}

func New(opts LoaderOptions) *loader {
    client := redis.NewClient(opts.Redis)
    return &loader{
        locations: []string{opts.FixtureDir},
        client:    client,
    }
}

func (l *loader) Load(names []string) error {
    ctx := parser.NewContext()
    fileParser := parser.New(l.locations)
    fixtureList, err := fileParser.ParseFiles(ctx, names)
    if err != nil {
        return err
    }
    return l.loadData(fixtureList)
}

func (l *loader) loadRedisDatabase(ctx context.Context, dbID int, db parser.Database, needTruncate bool) error {
    pipe := l.client.Pipeline()
    err := pipe.Select(ctx, dbID).Err()
    if err != nil {
        return err
    }

    if needTruncate {
        if err := pipe.FlushDB(ctx).Err(); err != nil {
            return err
        }
    }

    if db.Keys != nil {
        for k, v := range db.Keys.Values {
            if err := pipe.Set(ctx, k, v.Value, v.Expiration).Err(); err != nil {
                return err
            }
        }
    }

    if db.Sets != nil {
        for setKey, v := range db.Sets.Values {
            for v := range v.Values {
                if err := pipe.SAdd(ctx, setKey, v).Err(); err != nil {
                    return err
                }
            }
        }
    }

    if db.Maps != nil {
        for mapKey, v := range db.Maps.Values {
            for k, v := range v.Values {
                if err := pipe.HSet(ctx, mapKey, k, v).Err(); err != nil {
                    return err
                }
            }
        }
    }

    if _, err := pipe.Exec(ctx); err != nil {
        return err
    }

    return nil
}

func (l *loader) loadData(fixtures []*parser.Fixture) error {
    truncatedDatabases := make(map[int]struct{})

    for _, redisFixture := range fixtures {
        for dbID, db := range redisFixture.Databases {
            var needTruncate bool
            if _, ok := truncatedDatabases[dbID]; !ok {
                truncatedDatabases[dbID] = struct{}{}
                needTruncate = true
            }
            err := l.loadRedisDatabase(context.Background(), dbID, db, needTruncate)
            if err != nil {
                return err
            }
        }
    }
    return nil
}
