package utils

import (
	"github.com/shopspring/decimal"
	"math"
)

// 计算总价
func CalculateTotalPrice(price float32, num int32) float32 {
	value, _ := decimal.NewFromFloat32(price).Mul(decimal.NewFromInt32(num)).Round(2).Float64()
	return float32(value)
}

// 计算折扣
func CalculateRate(price1 float32, price2 float32) float32 {
	if price2 <= 0 {
		return 0
	}
	value, _ := decimal.NewFromFloat32(price1).Div(decimal.NewFromFloat32(price2)).Mul(decimal.NewFromInt32(100)).Round(2).Float64()
	return float32(value)
}

// 计算折扣价格
func CalculatePriceByRate(price float32, rate float32) float32 {
	value, _ := decimal.NewFromFloat32(price).Mul(decimal.NewFromFloat32(rate)).Div(decimal.NewFromFloat32(100)).Float64()
	return float32(value)
}

// 计算单位价格
func CalculateUnitPrice(price float32, num int32) float32 {
	if num <= 0 {
		return 0
	}
	value, _ := decimal.NewFromFloat32(price).Div(decimal.NewFromInt32(num)).Round(2).Float64()
	return float32(value)
}


// 获取两个坐标点距离(返回m)
func GetDistance(lng1, lat1, lng2, lat2 float64) float64 {
	radius := 6378137.0
	rad := math.Pi / 180.0
	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad
	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2)+math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta)) * radius
	dist, _ = decimal.NewFromFloat32(float32(dist)).Round(2).Float64()

	return dist
}

/*
 * 射线法判断坐标点是否在多边形内(经度对应Y轴，纬度对应X轴)
 * @param point map[string]float64 坐标点:x纬度，y经度
 * @param pList []map[string]float64 多边形坐标点集合
 * @return bool
 */
func IsPointInRegion(point map[string]float64, pList []map[string]float64) bool {
	var nCross int
	if pList == nil || len(pList) == 0 {
		return false
	}
	if _, ok := point["x"]; !ok {
		return false
	}
	if _, ok := point["y"]; !ok {
		return false
	}

	var cnt = len(pList)

	for i := 0; i < cnt; i++ {
		if _, ok := pList[i]["x"]; !ok {
			return false
		}
		if _, ok := pList[i]["y"]; !ok {
			return false
		}

		// 相邻的两个点
		var p1 = pList[i]
		var p2 = pList[(i+1)%cnt]
		if _, ok := p1["x"]; !ok {
			return false
		}
		if _, ok := p1["y"]; !ok {
			return false
		}
		if _, ok := p2["x"]; !ok {
			return false
		}
		if _, ok := p2["y"]; !ok {
			return false
		}

		var minPintX, maxPintX, minPintY, maxPintY = p1["x"], p2["x"], p1["y"], p2["y"]
		if p1["x"] > p2["x"] {
			minPintX = p2["x"]
			maxPintX = p1["x"]
		}
		if p1["y"] > p2["y"] {
			minPintY = p2["y"]
			maxPintY = p1["y"]
		}

		// 若此两相邻的点与X轴平行
		if p1["y"] == p2["y"] {
			if p1["y"] == point["y"] {
				// 三点共线，则判断点是否在中间,如果在中间，则必然在两点间，则射线（不是直线）必与此多边形有一交点
				if point["x"] >= minPintX && point["x"] <= maxPintX {
					nCross++
				}
			}

			//如果点point不在p1 p2之间，那此时过point的射线要么与p1 p2这条边有无数个交点 要么交点个数为0，此时不拿此边作为参考边
			continue
		}

		// 交点没在p1 p2这条边上，而在此条边的延长线上
		if point["y"] < minPintY || point["y"] >= maxPintY {
			continue
		}

		// 求出交点的坐标x
		var x = (point["y"]-p1["y"])*(p2["x"]-p1["x"])/(p2["y"]-p1["y"]) + p1["x"]
		// 统计左射线或右射线与边的交点都可以，此处统计的是右射线与多边形边的交点
		if x > point["x"] {
			nCross++
		}
	}

	if nCross%2 == 1 {
		return true
	} else {
		return false
	}
}

/*
 * 判断坐标点是否在圆内(包括圆上)(经度对应Y轴，纬度对应X轴)
 * @param point map[string]float64 坐标点:x纬度，y经度
 * @param cPoint map[string]float64 圆点坐标
 * @param radius float32 半径(km)
 * @return bool
 */
func IsPointInCircular(point map[string]float64, cPoint map[string]float64, radius float32) bool {
	if _, ok := point["x"]; !ok {
		return false
	}
	if _, ok := point["y"]; !ok {
		return false
	}
	if _, ok := cPoint["x"]; !ok {
		return false
	}
	if _, ok := cPoint["y"]; !ok {
		return false
	}
	if radius <= 0 {
		return false
	}

	var distance = float32(GetDistance(point["y"], point["x"], cPoint["y"], cPoint["x"]))
	if distance > radius*1000 {
		return false
	} else {
		return true
	}
}