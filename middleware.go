package main

import "strings"

func makeExecutionPoint() ExecutionPoint {
	return func(req BotRequest) (BotResponse, error) {
		resp := BotResponse{
			channelID: req.channelID,
		}

		addendum := " -d --account " + req.account

		output := Execute(req.message + addendum)
		if strings.Contains(output, sessionExpiredMessage) {
			authOutput := Execute(authCommand)
			Println(authOutput)
			Println()
			output = Execute(req.message + addendum)
		}

		resp.message = output

		return resp, nil
	}
}

func makeCheckAllowedCommandMiddleWare() Middleware {
	return func(next ExecutionPoint) ExecutionPoint {
		return func(req BotRequest) (BotResponse, error) {
			resp := BotResponse{
				channelID: req.channelID,
				message:   "",
			}

			if !req.approved && !IsAllowed(req.message) && !strings.HasSuffix(req.message, "help") {
				resp.message = "specified command is not allowed"
				return resp, nil
			}

			resp, err := next(req)

			return resp, err
		}
	}
}

var cmdsPendingApproval = map[string]string{}

func makeCheckApprovalRequiredCommandMiddleWare() Middleware {
	return func(next ExecutionPoint) ExecutionPoint {
		return func(req BotRequest) (BotResponse, error) {
			resp := BotResponse{
				channelID: req.channelID,
				message:   "",
			}

			if !req.approved && IsApprovalRequired(req.message) && !strings.HasSuffix(req.message, "help") {
				resp.message = "'" + req.message + "' requires approval. have your peers approve by saying @morgan pineapple."
				cmdsPendingApproval["pineapple"] = req.message
				return resp, nil
			}

			resp, err := next(req)

			return resp, err
		}
	}
}

func makeCheckForApprovalKeywordMiddleWare() Middleware {
	return func(next ExecutionPoint) ExecutionPoint {
		return func(req BotRequest) (BotResponse, error) {
			resp := BotResponse{
				channelID: req.channelID,
				message:   "",
			}

			if req.message == "pineapple" {
				req.message = cmdsPendingApproval["pineapple"]
				req.approved = true
				delete(cmdsPendingApproval, "pineapple")
				return next(req)
			}

			resp, err := next(req)

			return resp, err
		}
	}
}
