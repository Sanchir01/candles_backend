package userFeature

//func userLoaderFetch(keys []uuid.UUID) ([]*model.User, []error) {
//	// Получаем пользователей по их UUID
//	userStr := pgstoreuser.New(pgxdb)
//
//	users, err := userStr.GetById(context.Background(), keys)
//	if err != nil {
//		return nil, nil
//	}
//
//	// Подготовим результаты для DataLoader
//	results := make([]*model.User, len(users))
//	for i, user := range users {
//		results[i] = user
//	}
//
//	// Возвращаем пользователей и возможные ошибки
//	return results, nil
//}
