# Wallbox-Monitor

### Build the binaries & run the Ansible playbook:
```zsh
cd src && bash build.sh && cd .. && ansible-playbook -i inventory main.yml
# If podman container:
ansible-playbook -i inventory main.yml
```

### Example
```zsh
journalctl -u wallbox-monitor -f
```

### Build Container
```zsh
podman build -t wallbox-monitor .
podman run -d -p 0.0.0.0:2114:2114 wallbox-monitor
```

## Prometheus Integration

The application exposes metrics for Prometheus at `http://localhost:2114/metrics`.

### Wallbox metrics available:
- `wallbox_car_status` - Wallbox car connection status (0=unknown, 1=idle, 2=charging, 3=wait car, 4=finished)
- `wallbox_amp_limit_amps` - Requested Wallbox current limit in Amperes
- `wallbox_voltage_volts` - Measured Wallbox voltage in volts
- `wallbox_current_amps` - Measured Wallbox current in amps
- `wallbox_power_watts` - Wallbox power in watts
- `wallbox_total_energy_wh` - Wallbox total charged energy in watt-hours
- `wallbox_session_energy_dws` - Wallbox session energy in deciwatt-seconds
- `wallbox_phase_mode` - Wallbox supply phase mode (1=force single-phase, 2=force three-phase)
- `wallbox_override_state` - Wallbox charge control override state (0=automatic, 1=force off, 2=force on)
- `wallbox_charge_allowed` - Whether Wallbox is allowed to charge (1=allowed, 0=disallowed)
- `wallbox_temp_celsius{sensor="0"}` - Wallbox internal temperature sensor 0 in degrees Celsius
- `wallbox_temp_celsius{sensor="1"}` - Wallbox internal temperature sensor 1 in degrees Celsius (if present)

### Prometheus Configuration:
```yaml
scrape_configs:
  - job_name: 'wallbox-monitor'
    static_configs:
      - targets: ['localhost:2114']
    scrape_interval: 5s
```

## Build container

```zsh
mkdir -p /data/solar-monitor ; chmod 1777 /data/solar-monitor
TMPDIR=/tmp docker build -t solar-monitor:latest .
TMPDIR=/tmp docker run -d --name solar-monitor -p 2112:2112 -v /data/solar-monitor:/home/appuser/solar-monitor:Z --restart unless-stopped solar-monitor:latest

# OR run in foreground
TMPDIR=/tmp docker run -p 2112:2112 -v /data/solar-monitor:/home/appuser/solar-monitor:Z --rm solar-monitor:latest
```

---

## ⚙️ Configuration Parameter Reference

The monitoring engine is completely dynamic and reads its environment configurations on startup. If an environment variable is missing, it safely falls back to its built-in out-of-the-box defaults.

### Core Service & Runtime Variables

| Environment Variable | Description | Expected Type | Default Value |
| :--- | :--- | :--- | :--- |
| **`WALLBOX_IP`** | Network destination IP address of your physical go-eCharger hardware. | String | `192.168.30.40` |
| **`WALLBOX_PORT`** | Active HTTP port configuration of the physical wallbox hardware API. | String | `80` |
| **`PROMETHEUS_PORT`** | Active network port binding configuration for exposing Prometheus metrics. | String | `2114` |
| **`FORCE_IPV4`** | Binds network socket connection layers strictly to the IPv4 stack. | Boolean | `true` |
| **`OUTPUT_JSON`** | Toggles whether compiled metrics are cleanly printed as JSON strings to stdout. | Boolean | `true` |
| **`WRITE_JSON_TO_FILE`** | Toggles whether live state updates are continuously written to a local JSON file. | Boolean | `true` |
| **`ENABLE_PROMETHEUS`** | Controls execution of the internal Prometheus metric scraping web server endpoint. | Boolean | `true` |
| **`ENABLE_HOME_ASSISTANT`** | Master switch to toggle outbound state synchronization updates down to Home Assistant. | Boolean | `false` |

---

### Home Assistant Integration Variables

| Environment Variable | Description | Expected Type | Default Value |
| :--- | :--- | :--- | :--- |
| **`URL_HOME_ASSISTANT`** | Target API URL endpoint of your active Home Assistant server instance. | String | `http://192.168.1.8:8123` |
| **`TOKEN_HOME_ASSISTANT`** | Long-Lived Access Token generated from Home Assistant for secure REST authentication. | String | *(Uses dynamic secure token)* |
| **`API_BOOL_WALLBOX_CHARGE_ALLOWED`** | Target entity API path mapping for checking if charging authorization is clear. | String | `/api/states/input_boolean.wallbox_charge_allowed` |
| **`API_INT_WALLBOX_AMP_LIMIT`** | Target entity API path mapping for tracking requested charging limit parameters in Amperes. | String | `/api/states/input_number.wallbox_amp_limit_amps` |
| **`API_INT_WALLBOX_POWER_W`** | Target entity API path mapping for tracking real-time electrical wallbox draw in Watts. | String | `/api/states/input_number.wallbox_power_w` |
| **`API_INT_WALLBOX_TOTAL_ENERGY_WH`** | Target entity API path mapping for lifetime energy delivery accumulated metrics (Watt-hours). | String | `/api/states/input_number.wallbox_total_energy_wh` |
| **`API_INT_WALLBOX_SESSION_ENERGY_DWS`** | Target entity API path mapping for the active charging loop session energy (Deciwatt-seconds). | String | `/api/states/input_number.wallbox_session_energy_dws` |
| **`API_INT_WALLBOX_PHASE_MODE`** | Target entity API path mapping for tracked grid phase distribution (Single-Phase vs Three-Phase). | String | `/api/states/input_number.wallbox_phase_mode` |
| **`API_INT_WALLBOX_OVERRIDE_STATE`** | Target entity API path mapping for manual override status tracking states. | String | `/api/states/input_number.wallbox_override_state` |


---

# WALLBOX API (go-e Charger Gemini flex 11kW)

```
http://192.168.30.40/api/status # Fetches a single JSON object containing every single live state, configuration parameter, and sensor measurement.
http://192.168.30.40/api/status?filter=all # Alternate explicit call ensuring all monitoring keys are generated uncompressed by the device.
http://192.168.30.40/api/status?filter=car # Retrieves only the car status key (0: unknown, 1: idle, 2: charging, 3: wait car, 4: finished).
http://192.168.30.40/api/status?filter=amp # Retrieves only the requested current limit setting in whole Amperes.
http://192.168.30.40/api/status?filter=nrg # Retrieves only the electrical measurement array containing voltage, current, and total power metrics.
http://192.168.30.40/api/status?filter=wh # Retrieves only the life-to-date total charged energy recorded by the hardware internal meter in Watt-hours.
http://192.168.30.40/api/status?filter=dws # Retrieves only the energy delivered during the current charging session in Deciwatt-seconds.
http://192.168.30.40/api/status?filter=psm # Retrieves only the current supply phase mode (1: force single-phase, 2: force three-phase).
http://192.168.30.40/api/status?filter=frc # Retrieves only the charge control override state (0: automatic, 1: force off, 2: force on).
http://192.168.30.40/api/status?filter=tmp # Retrieves only the device internal microprocessor temperatures array in degrees Celsius.
http://192.168.30.40/api/status?filter=alw # Retrieves only the clear-to-charge status flag showing if vehicle power delivery is allowed.
http://192.168.30.40/api/status?filter=car,amp,nrg,wh,dws,psm,frc,tmp,alw # Custom multi-filter call retrieving all key metrics simultaneously in one request.
```