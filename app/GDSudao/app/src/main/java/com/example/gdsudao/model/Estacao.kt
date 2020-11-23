package com.example.gdsudao.model

data class Estacao (
    var _id: String?,
    var altitude: Float?,
    var codigoCPTEC: String?,
    var codigoINMET: String?,
    var localizacao: Localizacao?,
    var nomeEstacao: String?,
    var tipoEstacao: String?,
    var uf: String?
)

