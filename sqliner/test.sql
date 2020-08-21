SELECT  employer_id AS id,
        t.value AS region,
        e.name AS employer_name,
        'https://hh.ru/admin/employer.do?employerId=' || employer_id AS employer_page,
        concat_ws(' ', u.last_name_cache, u.first_name_cache, u.middle_name_cache ) AS manager_name,
        'Есть доступ в отложенных' AS reason
FROM employer AS e 
    JOIN translation AS t ON 'area.' || e.area_id = t.name
    JOIN hhuser AS u ON e.manager_id = u.user_id
WHERE e.employer_id IN (
    SELECT a.employer_id AS employer_id FROM account AS a 
        JOIN cart AS c ON a.account_id = c.account_id
        JOIN services_cart AS sc ON c.cart_id = sc.cart_id 
        JOIN service AS s ON sc.service_id = s.service_id 
        WHERE s.service_type IN ('FA', 'FA+VPPL', 'FA+VPP', 'CIV+VPPL') 
            AND c.status IN (1, 2)
            AND a.seller_account_id = 1
    )
    AND t.lang = 'RU';
