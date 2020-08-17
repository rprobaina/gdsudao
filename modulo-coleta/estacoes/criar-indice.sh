# Cria indice
db.estacoes.createIndex({ "localizacao.coordenadas" : "2dsphere"})
 
# {
# 	"createdCollectionAutomatically" : false,
# 	"numIndexesBefore" : 1,
# 	"numIndexesAfter" : 2,
# 	"ok" : 1
# }
#