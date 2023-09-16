package debugDto

import (
	"fmt"
	"ilserver/domain"
)

/* example:

{
  "some": [
    [
      "<topic_name>",
      <lang>
    ]
  ]
}

*/

type Topic []interface{}

type Topics struct {
	Some []Topic `json:"some"`
}

func TopicsToDomain(topics Topics) (error, []domain.Topic) {
	var result []domain.Topic
	for _, topic := range topics.Some {
		if len(topic) != 2 {
			return fmt.Errorf("topic as list has incorrect size"), nil
		}

		// ***

		lang, convOk := topic[1].(float64)
		if !convOk {
			return fmt.Errorf("lang has incorrect type"), nil
		}

		topicName, convOk := topic[0].(string)
		if !convOk {
			return fmt.Errorf("topic name has incorrect type"), nil
		}

		// ***

		var one domain.Topic
		one.Lang = int(lang)
		one.Name = topicName
		result = append(result, one)
	}
	return nil, result
}
