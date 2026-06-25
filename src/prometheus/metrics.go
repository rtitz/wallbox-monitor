package prometheus

import (
	"strconv"

	"wallbox-monitor/wallboxApi"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	WallboxCarStatus = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wallbox_car_status",
		Help: "Wallbox car connection status (0=unknown, 1=idle, 2=charging, 3=wait car, 4=finished)",
	})

	WallboxAmpLimitAmps = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wallbox_amp_limit_amps",
		Help: "Requested Wallbox current limit in Amperes",
	})

	WallboxVoltageVolts = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wallbox_voltage_volts",
		Help: "Wallbox voltage in volts",
	})

	WallboxCurrentAmps = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wallbox_current_amps",
		Help: "Wallbox current in amps",
	})

	WallboxPowerWatts = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wallbox_power_watts",
		Help: "Wallbox power in watts",
	})

	WallboxTotalEnergyWh = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wallbox_total_energy_wh",
		Help: "Wallbox total charged energy in watt-hours",
	})

	WallboxSessionEnergyDws = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wallbox_session_energy_dws",
		Help: "Wallbox session energy in deciwatt-seconds",
	})

	WallboxPhaseMode = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wallbox_phase_mode",
		Help: "Wallbox supply phase mode (1=force single-phase, 2=force three-phase)",
	})

	WallboxOverrideState = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wallbox_override_state",
		Help: "Wallbox charge control override state (0=automatic, 1=force off, 2=force on)",
	})

	WallboxChargeAllowed = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wallbox_charge_allowed",
		Help: "Whether Wallbox is allowed to charge (1 = allowed, 0 = disallowed)",
	})

	WallboxTempCelsius = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "wallbox_temp_celsius",
		Help: "Wallbox internal temperature sensors in degrees Celsius",
	}, []string{"sensor"})
)

func UpdateMetrics(status *wallboxApi.Status) {
	WallboxCarStatus.Set(float64(status.Car))
	WallboxAmpLimitAmps.Set(float64(status.Amp))

	voltage, current, power := 0.0, 0.0, 0.0
	if len(status.Nrg) > 0 {
		// nrg[0] is Phase 1 voltage in volts
		voltage = status.Nrg[0]
	}
	if len(status.Nrg) > 4 {
		// nrg[4..6] are phase currents in 0.1 A
		current += status.Nrg[4] / 10.0
	}
	if len(status.Nrg) > 5 {
		current += status.Nrg[5] / 10.0
	}
	if len(status.Nrg) > 6 {
		current += status.Nrg[6] / 10.0
	}
	if len(status.Nrg) > 11 {
		// nrg[11] is total power in 0.1 W
		power = status.Nrg[11] / 10.0
	}

	WallboxVoltageVolts.Set(voltage)
	WallboxCurrentAmps.Set(current)
	WallboxPowerWatts.Set(power)
	WallboxTotalEnergyWh.Set(status.Wh)
	WallboxSessionEnergyDws.Set(status.Dws)
	WallboxPhaseMode.Set(float64(status.Psm))
	WallboxOverrideState.Set(float64(status.Frc))
	if status.Alw {
		WallboxChargeAllowed.Set(1)
	} else {
		WallboxChargeAllowed.Set(0)
	}

	for idx, tma := range status.Tma {
		WallboxTempCelsius.WithLabelValues(strconv.Itoa(idx)).Set(tma)
	}
}
