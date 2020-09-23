package repo

const (
	stmtAllDevice = `
	SELECT a.device_token AS device_token,
		a.device_type AS device_type
	FROM analytic.ios_device_token AS a 
		LEFT JOIN analytic.ios_invalid_token AS i
		ON UNHEX(a.device_token) = i.device_token
		INNER JOIN analytic.ios_zone_offset AS zone
		ON a.timezone = zone.timezone
	WHERE (i.device_token IS NULL) 
		AND (
			HOUR(
				DATE_ADD(UTC_TIMESTAMP(), INTERVAL zone.utc_offset HOUR)
			) BETWEEN 6 AND 24
		);`

	stmtTestDevice = `
	SELECT LOWER(HEX(token)) AS device_token, 
		'phone' AS device_type
	FROM analytic.ios_test_device;`

	selectPaidDevice = stmtTestDevice
)
