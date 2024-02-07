package routes

import (
	"shorten_url/database"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// 將舊連接儲存在db，當有人使用新連結會去和db比對來找到原本網址

func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")
	r := database.CreateClient(0)
	defer r.Close()

	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short not found in the database"})
	} else if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "cannot connect to RDB"})
	}

	rInr := database.CreateClient(1)
	defer rInr.Close()

	_ = rInr.Incr(database.Ctx, "counter")
	return c.Redirect(value, 301)
}
