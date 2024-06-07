#!/bin/bash

INTERVAL="1"

main() {
	while true; do
		output_download=""
		output_upload=""
		output_download_unit=""
		output_upload_unit=""

		initial_download=$(cat /sys/class/net/eth0/statistics/rx_bytes)
		initial_upload=$(cat /sys/class/net/eth0/statistics/tx_bytes)

		sleep "$INTERVAL"

		final_download=$(cat /sys/class/net/eth0/statistics/rx_bytes)
		final_upload=$(cat /sys/class/net/eth0/statistics/tx_bytes)

		total_download_bps=$(expr "$final_download" - "$initial_download")
		total_upload_bps=$(expr "$final_upload" - "$initial_upload")

		if [ "$total_download_bps" -gt 1073741824 ]; then
			output_download=$(echo "$total_download_bps 1024" | awk '{printf "%.1f \n", $1/($2 * $2 * $2)}')
			output_download_unit="gb"
		elif [ "$total_download_bps" -gt 1048576 ]; then
			output_download=$(echo "$total_download_bps 1024" | awk '{printf "%.1f \n", $1/($2 * $2)}')
			output_download_unit="mb"
		else
			output_download=$(echo "$total_download_bps 1024" | awk '{printf "%.1f \n", $1/$2}')
			output_download_unit="kb"
		fi

		if [ "$total_upload_bps" -gt 1073741824 ]; then
			output_upload=$(echo "$total_download_bps 1024" | awk '{printf "%.1f \n", $1/($2 * $2 * $2)}')
			output_upload_unit="gb"
		elif [ "$total_upload_bps" -gt 1048576 ]; then
			output_upload=$(echo "$total_upload_bps 1024" | awk '{printf "%.1f \n", $1/($2 * $2)}')
			output_upload_unit="mb"
		else
			output_upload=$(echo "$total_upload_bps 1024" | awk '{printf "%.1f \n", $1/$2}')
			output_upload_unit="kb"
		fi

		echo "$output_upload$output_upload_unit|$output_download$output_download_unit" >~/scripts/custom_output/net_speed

	done
}
main
