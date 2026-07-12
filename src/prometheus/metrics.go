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

	WallboxUserToken = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "wallbox_user_token",
		Help: "The active RFID user token ID (0=Open, 1=RFID 1, 2=Guest, etc.)",
	})
)

func UpdateMetrics(status *wallboxApi.Status) {
	WallboxCarStatus.Set(float64(status.Car))
	WallboxAmpLimitAmps.Set(float64(status.Amp))

	voltage, current, power := 0.0, 0.0, 0.0
	nrgLen := len(status.Nrg)

	// nrg[0] is Phase 1 voltage in volts
	if nrgLen > 0 {
		voltage = status.Nrg[0]
	}

	// nrg[4..6] are phase currents in Amperes for API v2
	if nrgLen > 4 {
		current += status.Nrg[4]
	}
	if nrgLen > 5 {
		current += status.Nrg[5]
	}
	if nrgLen > 6 {
		current += status.Nrg[6]
	}

	// nrg[11] is total power in Watts for API v2
	if nrgLen > 11 {
		power = status.Nrg[11]
	}

	WallboxVoltageVolts.Set(voltage)
	WallboxCurrentAmps.Set(current)
	WallboxPhaseMode.Set(float64(status.Psm))
	WallboxOverrideState.Set(float64(status.Frc))

	// ROUTING LOGIC: Isolate your car (Ust 0 or 1) from any Guests (Ust 2+)
	if status.Ust == 0 || status.Ust == 1 {
		// Your car is charging: Update existing metrics normally
		WallboxPowerWatts.Set(power)
		WallboxTotalEnergyWh.Set(status.Wh)
		WallboxSessionEnergyDws.Set(status.Dws)
	} else {
		// A Guest is charging: Freeze your car metrics
		WallboxPowerWatts.Set(0) // Set power to 0 so your solar/grid rules ignore the guest
		// Do NOT update WallboxTotalEnergyWh and WallboxSessionEnergyDws here.
		// This keeps your car odometer frozen at its last value while the guest charges.
	}

	if status.Alw {
		WallboxChargeAllowed.Set(1)
	} else {
		WallboxChargeAllowed.Set(0)
	}

	for idx, tma := range status.Tma {
		WallboxTempCelsius.WithLabelValues(strconv.Itoa(idx)).Set(tma)
	}

	WallboxUserToken.Set(float64(status.Ust))
}
