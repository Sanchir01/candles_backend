package response

import "github.com/Sanchir01/sandjma_graphql/internal/gql/model"

func NewInternalErrorProblem(message string) model.InternalErrorProblem {
	return model.InternalErrorProblem{Message: message}
}

func NewVersionMismatchProblem() model.VersionMismatchProblem {
	return model.VersionMismatchProblem{Message: "version mismatch"}
}
