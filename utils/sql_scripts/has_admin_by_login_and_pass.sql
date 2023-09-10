SELECT count(*) AS RecordCount 
FROM [Admins] 
WHERE [Admins].[Login] = 'admin' 
AND [Admins].[Pass] =  'new_pass';