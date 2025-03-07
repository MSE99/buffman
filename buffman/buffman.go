package buffman

import "context"

func StartDispatchToFMA(ctx context.Context) error {
	tk, err := newFmaToken(ctx)
	if err != nil {
		return err
	}
	go tk.waitAndRefresh()

	return nil
}
