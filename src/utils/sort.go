package utils

import (
	"math/rand"
	"src/structs"
)

func sortFullCommentsByDate(comments []structs.FullComment) []structs.FullComment {
	if len(comments) < 2 {
		return comments
	}

	left, right := 0, len(comments)-1
	pivot := rand.Intn(len(comments))
	comments[pivot], comments[right] = comments[right], comments[pivot]

	for i := 0; i < len(comments); i++ {
		if comments[i].PostedTime.After(comments[right].PostedTime) {
			comments[left], comments[i] = comments[i], comments[left]
			left++
		}
	}

	comments[left], comments[right] = comments[right], comments[left]

	sortFullCommentsByDate(comments[:left])
	sortFullCommentsByDate(comments[left+1:])

	return comments
}

func CommentSort(array []structs.FullComment) []structs.FullComment {
	var result []structs.FullComment
	primaryCommentsCounter := 0

	for i := 0; i < len(array); i++ {
		if array[i].ParentID == 0 {
			result = append(result, array[i])
			primaryCommentsCounter++
		}
	}
	result = sortFullCommentsByDate(result)
	resultCPY := result
	index := 1
	iCounter := 0

	for i := 0; i < len(resultCPY); i++ {
		var tempArray []structs.FullComment
		secondaryCommentsCounter := 0
		for j := 0; j < len(array); j++ {
			if array[j].ParentID == resultCPY[i].ID {
				tempArray = append(tempArray, array[j])
				secondaryCommentsCounter++
			}
		}

		tempArray = sortFullCommentsByDate(tempArray)
		for j := 0; j < len(tempArray); j++ {
			result = insertInArray(result, tempArray[j], index)
			index++

		}
		iCounter++
		index++
	}

	return result
}

func SortHomePageReturnByLikesCroissant(posts []structs.HomePageReturn) []structs.HomePageReturn {
	if len(posts) < 2 {
		return posts
	}

	left, right := 0, len(posts)-1
	pivot := rand.Intn(len(posts))
	posts[pivot], posts[right] = posts[right], posts[pivot]

	for i := 0; i < len(posts); i++ {
		if posts[i].Likes < posts[right].Likes {
			posts[left], posts[i] = posts[i], posts[left]
			left++
		}
	}

	posts[left], posts[right] = posts[right], posts[left]

	SortHomePageReturnByLikesCroissant(posts[:left])
	SortHomePageReturnByLikesCroissant(posts[left+1:])

	return posts
}

func SortHomePageReturnByLikesDecroissant(posts []structs.HomePageReturn) []structs.HomePageReturn {
	if len(posts) < 2 {
		return posts
	}

	left, right := 0, len(posts)-1
	pivot := rand.Intn(len(posts))
	posts[pivot], posts[right] = posts[right], posts[pivot]

	for i := 0; i < len(posts); i++ {
		if posts[i].Likes > posts[right].Likes {
			posts[left], posts[i] = posts[i], posts[left]
			left++
		}
	}

	posts[left], posts[right] = posts[right], posts[left]

	SortHomePageReturnByLikesDecroissant(posts[:left])
	SortHomePageReturnByLikesDecroissant(posts[left+1:])

	return posts
}

func SortHomePageReturnByDateDecroissant(posts []structs.HomePageReturn) []structs.HomePageReturn {
	var result []structs.HomePageReturn

	for i := len(posts) - 1; i >= 0; i-- {
		result = append(result, posts[i])
	}

	return posts
}
