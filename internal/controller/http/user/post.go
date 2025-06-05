package user

// POST /users
// func (c *Controller) CreateUser(eCtx echo.Context) error {
// 	var newUser model.User
// 	if err := eCtx.Bind(&newUser); err != nil {
// 		return eCtx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
// 	}

// 	newUser.ID = len(users) + 1
// 	users = append(users, newUser)

// 	return eCtx.JSON(http.StatusCreated, newUser)
// }
