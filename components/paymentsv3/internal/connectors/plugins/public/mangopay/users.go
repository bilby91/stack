package mangopay

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/paymentsv3/internal/models"
)

type usersState struct {
	LastPage         int       `json:"last_page"`
	LastCreationDate time.Time `json:"last_creation_date"`
}

func (p Plugin) fetchNextUsers(ctx context.Context, req models.FetchNextOthersRequest) (models.FetchNextOthersResponse, error) {
	var oldState usersState
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchNextOthersResponse{}, err
		}
	}

	newState := usersState{
		LastPage:         oldState.LastPage,
		LastCreationDate: oldState.LastCreationDate,
	}

	var others []json.RawMessage
	hasMore := false
	for page := oldState.LastPage; ; page++ {
		newState.LastPage = page

		pagedUsers, err := p.client.GetUsers(ctx, page, req.PageSize)
		if err != nil {
			return models.FetchNextOthersResponse{}, err
		}

		if len(pagedUsers) == 0 {
			break
		}

		for _, user := range pagedUsers {
			userCreationDate := time.Unix(user.CreationDate, 0)
			switch userCreationDate.Compare(oldState.LastCreationDate) {
			case -1, 0:
				// creationDate <= state.LastCreationDate, nothing to do,
				// we already processed this user.
				continue
			default:
			}

			raw, err := json.Marshal(user)
			if err != nil {
				return models.FetchNextOthersResponse{}, err
			}

			others = append(others, raw)

			newState.LastCreationDate = userCreationDate

			if len(others) >= req.PageSize {
				break
			}
		}

		if len(pagedUsers) < req.PageSize {
			break
		}

		if len(others) >= req.PageSize {
			hasMore = true
			break
		}
	}

	payload, err := json.Marshal(newState)
	if err != nil {
		return models.FetchNextOthersResponse{}, err
	}

	return models.FetchNextOthersResponse{
		Others:   others,
		NewState: payload,
		HasMore:  hasMore,
	}, nil
}
