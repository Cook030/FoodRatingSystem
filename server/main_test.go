package main

import (
	"foodRatingSystem/database"
	"foodRatingSystem/model"
	"foodRatingSystem/service"
	"testing"
)

func TestSubmitReviewIntegration(t *testing.T) {
	database.Connectdb()

	targetResName := "南湖美食广场"
	testUserID := "test_user_999"
	testUserName := "test_user_name_999"
	testStars := 5.0
	testComment := "单元测试自动生成的评价"

	var oldScore float64
	var oldRefID uint
	var r model.Restaurant
	err := database.DB.Where("name = ?", targetResName).First(&r).Error
	if err != nil {
		t.Fatalf("无法获取测试餐厅数据: %v", err)
	}
	oldRefID = r.ID
	oldScore = r.AverageScore

	err = service.SubmitReview(targetResName, testUserID, testUserName, testStars, testComment)
	if err != nil {
		t.Errorf("SubmitReview 执行报错: %v", err)
	}

	var newScore float64
	var newCount int64
	var rest model.Restaurant
	database.DB.Model(&model.Restaurant{}).Where("id = ?", oldRefID).Scan(&rest)
	newScore = rest.AverageScore
	database.DB.Model(&model.Rating{}).Where("restaurant_id = ?", oldRefID).Count(&newCount)

	if newCount <= 0 {
		t.Error("校验失败：评价总数没有增加")
	}

	t.Logf("✅ 测试通过！餐厅ID: %d, 旧分数: %.2f, 新分数: %.2f, 总评价数: %d", oldRefID, oldScore, newScore, newCount)

	defer func() {
		database.DB.Where("user_id = ?", testUserID).Delete(&model.Rating{})
		t.Log("🧹 测试数据已清理")
	}()
}
