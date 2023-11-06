# Battery adapter Modbus Interface

| Address | Type | Function Code | Description                    | Unit         |
|---------------------|-----------------|-------------------|--------------------------------|---------------|
| 0                   | uint16          | Input Register | soc          |               |
| 1                   | uint16          | Input Register | soh           |               |
| 2                   | float32          | Input Register | capacity_Wh   | Wh   |
| 4                   | float32          | Holding Register | Active Power  |       W        |
| 0                   | float32          | Holding Register | Setpoint  |       W        |


