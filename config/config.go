package config

type Config struct {
    Server *ServerConfig `json:"server"`
    Price  *PriceConfig  `json:"price"`
    DB     *DBConfig     `json:"db"`
    Nats   *NatsConfig   `json:"nats"`
    Poets  []*PoetConfig `json:"poets"`
}

type PriceConfig struct {
    Provider    string `json:"provider"`
    RefreshTime int    `json:"refreshTime"`
}

type ServerConfig struct {
    Port string `json:"port"`
}

type NatsConfig struct {
    Enabled bool   `json:"enabled"`
    Uri     string `json:"uri"`
}

type DBConfig struct {
    Uri string `json:"uri"`
}

type PoetConfig struct {
    Name     string        `json:"name"`
    Info     *PoetInfo     `json:"info"`
    Settings *PoetSettings `json:"settings"`
}

type PoetInfo struct {
    Description string `json:"description"`
    DiscordLink string `json:"discord-link"`
}

type PoetSettings struct {
    PhaseShift int `json:"phase-shift"`
    CycleGap   int `json:"cycle-gap"`
}
