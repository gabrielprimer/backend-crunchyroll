seleciona os content_source de um anime

SELECT cs.*
FROM content_sources cs
JOIN animes a ON cs.anime_id = a.id
WHERE a.slug = 'solo-leveling';





ver se eu consigo pegar o 
title, dub ou legendado, tempo do episodio, sinopsis do ep, data de lancamento





DAR COMMIT NA BRANCH DEV !!!!!!!!!!!!!!