package migrations

import "context"

type Migration struct {
	ID string
	Up func(ctx context.Context) error
}

func Run(ctx context.Context, migrations []Migration) error {
	for _, m := range migrations {
		if err := m.Up(ctx); err != nil {
			return err
		}
	}
	return nil
}