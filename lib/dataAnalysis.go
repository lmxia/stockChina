package lib

import (
	types "stockChina/type"
)

func ConventInfoToPoint(infos []types.StockPriceHistoryInfo) []types.CoordinatePoint {
	var points []types.CoordinatePoint
	i := 0.0
	for _, info := range infos {
		var point types.CoordinatePoint
		i = i + 1
		point.XPoint = i
		point.YPoint = info.ClosingPriceYesterday
		points = append(points, point)
	}
	return points
}

func GetDerivation(befores []types.CoordinatePoint) []types.CoordinatePoint {
	var afters []types.CoordinatePoint
	i := 0
	for i < len(befores) {
		i = i + 1
		var coordinatepoint types.CoordinatePoint
		coordinatepoint.XPoint = befores[i].XPoint
		coordinatepoint.YPoint = (befores[i].YPoint - befores[i-1].YPoint) / (befores[i].XPoint - befores[i-1].XPoint)
		afters = append(afters, coordinatepoint)
	}
	return afters
}
