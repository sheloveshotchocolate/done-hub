package controller

import (
	"done-hub/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type StatisticsByPeriod struct {
	UserStatistics       []*model.UserStatisticsByPeriod    `json:"user_statistics"`
	ChannelStatistics    []*model.LogStatisticGroupChannel  `json:"channel_statistics"`
	RedemptionStatistics []*model.RedemptionStatisticsGroup `json:"redemption_statistics"`
	OrderStatistics      []*model.OrderStatisticsGroup      `json:"order_statistics"`
}

func GetStatisticsByPeriod(c *gin.Context) {
	startTimestamp, _ := strconv.ParseInt(c.Query("start_timestamp"), 10, 64)
	endTimestamp, _ := strconv.ParseInt(c.Query("end_timestamp"), 10, 64)
	groupType := c.Query("group_type")
	userID, _ := strconv.Atoi(c.Query("user_id"))

	statisticsByPeriod := &StatisticsByPeriod{}

	userStatistics, err := model.GetUserStatisticsByPeriod(startTimestamp, endTimestamp)
	if err == nil {
		statisticsByPeriod.UserStatistics = userStatistics
	}

	startTime := time.Unix(startTimestamp, 0)
	endTime := time.Unix(endTimestamp, 0)
	startDate := startTime.Format("2006-01-02")
	endDate := endTime.Format("2006-01-02")
	channelStatistics, err := model.GetChannelExpensesStatisticsByPeriod(startDate, endDate, groupType, userID)

	if err == nil {
		statisticsByPeriod.ChannelStatistics = channelStatistics
	}

	redemptionStatistics, err := model.GetStatisticsRedemptionByPeriod(startTimestamp, endTimestamp)
	if err == nil {
		statisticsByPeriod.RedemptionStatistics = redemptionStatistics
	}

	orderStatistics, err := model.GetStatisticsOrderByPeriod(startTimestamp, endTimestamp)
	if err == nil {
		statisticsByPeriod.OrderStatistics = orderStatistics
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    statisticsByPeriod,
	})
}

type StatisticsDetail struct {
	UserStatistics      *model.StatisticsUser         `json:"user_statistics"`
	ChannelStatistics   []*model.ChannelStatistics    `json:"channel_statistics"`
	RedemptionStatistic []*model.RedemptionStatistics `json:"redemption_statistic"`
	OrderStatistics     []*model.OrderStatistics      `json:"order_statistics"`
	RpmTpmStatistics    *RpmTpmStatistics             `json:"rpm_tpm_statistics"`
}

type RpmTpmStatistics struct {
	RPM int64 `json:"rpm"`
	TPM int64 `json:"tpm"`
}

func GetStatisticsDetail(c *gin.Context) {

	statisticsDetail := &StatisticsDetail{}
	userStatistics, err := model.GetStatisticsUser()
	if err == nil {
		statisticsDetail.UserStatistics = userStatistics
	}

	channelStatistics, err := model.GetStatisticsChannel()
	if err == nil {
		statisticsDetail.ChannelStatistics = channelStatistics
	}

	redemptionStatistics, err := model.GetStatisticsRedemption()
	if err == nil {
		statisticsDetail.RedemptionStatistic = redemptionStatistics
	}

	orderStatistics, err := model.GetStatisticsOrder()
	if err == nil {
		statisticsDetail.OrderStatistics = orderStatistics
	}

	// 获取最近60秒的RPM和TPM统计
	rpmTpmStats, err := model.GetRpmTpmStatistics()
	if err == nil {
		statisticsDetail.RpmTpmStatistics = &RpmTpmStatistics{
			RPM: rpmTpmStats.RPM,
			TPM: rpmTpmStats.TPM,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    statisticsDetail,
	})
}
