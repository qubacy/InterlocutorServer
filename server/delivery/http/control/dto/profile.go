package dto

import (
	domain "ilserver/domain/memory"
)

// direct mapping
// -----------------------------------------------------------------------

type Profile struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Contact  string `json:"contact"`
}

func MakeProfile(profile domain.Profile) Profile {
	return Profile{
		Id:       profile.Id,
		Username: profile.Username,
		Contact:  profile.Contact,
	}
}

func MakeProfiles(profiles domain.ProfileList) []Profile {
	result := []Profile{}
	for i := range profiles {
		result = append(result,
			MakeProfile(profiles[i]))
	}
	return result
}
