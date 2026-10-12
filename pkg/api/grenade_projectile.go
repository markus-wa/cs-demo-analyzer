package api

import (
	"github.com/golang/geo/r3"

	"github.com/markus-wa/demoinfocs-golang/v6/pkg/demoinfocs/common"
)

// findGrenadeProjectileID returns the unique ID of the live projectile thrown with the given grenade.
// A player may have thrown several projectiles with the same grenade (weapon instances are re-used
// across throws), in that case the closest one to the event position is the right one.
func findGrenadeProjectileID(analyzer *Analyzer, grenadeUniqueID common.UniqueID, eventPosition r3.Vector) int64 {
	var projectileID int64
	bestDistance := -1.0
	bestEntityID := -1
	for _, projectile := range analyzer.parser.GameState().GrenadeProjectiles() {
		if projectile.WeaponInstance.UniqueID() != grenadeUniqueID {
			continue
		}

		distance := projectile.Position().Sub(eventPosition).Norm2()
		entityID := projectile.Entity.ID()
		isCloser := bestDistance == -1 || distance < bestDistance
		// On an exact distance tie either projectile is an equally valid match, so prefer the lower entity ID
		// to keep the result deterministic regardless of iteration order.
		isSameDistanceWithLowerEntityID := distance == bestDistance && entityID < bestEntityID
		if isCloser || isSameDistanceWithLowerEntityID {
			bestDistance = distance
			bestEntityID = entityID
			projectileID = int64(projectile.UniqueID())
		}
	}

	return projectileID
}
