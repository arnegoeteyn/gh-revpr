package faker

import (
	"fmt"

	"github.com/go-faker/faker/v4"
)

func ReviewWorktreeName() string {
	return fmt.Sprintf("review-%s-%s", faker.Word(), faker.Word())

}
