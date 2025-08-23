import React from "react";

function floatToPercent(value) {
  // Clamp value between 0.7 and 1.9
  const min = 0.7, max = 1.9;
  const clamped = Math.max(min, Math.min(max, value));
  // Convert to percentage
  return Math.round(((clamped - min) / (max - min)) * 100);
}

export default function BatteryIndicator({ value }) {
  const percent = floatToPercent(value);
  // Fill width for battery icon
  const fillWidth = `${percent}%`;

  return (
    <div style={{ display: "flex", alignItems: "center", gap: "8px" }}>
      <div style={{
        width: 40, height: 20, border: "2px solid #333", borderRadius: 4, position: "relative", background: "#eee"
      }}>
        <div style={{
          width: fillWidth, height: "100%", background: percent > 20 ? "#4caf50" : "#f44336", borderRadius: 2, transition: "width 0.3s"
        }} />
        <div style={{
          position: "absolute", right: -6, top: 6, width: 6, height: 8, background: "#333", borderRadius: 2
        }} />
      </div>
      <span style={{ fontWeight: "bold" }}>{percent}%</span>
    </div>
  );
}