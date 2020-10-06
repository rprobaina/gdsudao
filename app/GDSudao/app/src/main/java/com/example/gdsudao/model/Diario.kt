package com.example.gdsudao.model

data class Diario(
    var _id: String,
    var codigoCPTEC: String,
    var codigoINMET: String,
    var dataMedicao: String,
    var temperaturaMaxima: Float,
    var temperaturaMedia: Float,
    var temperaturaMinima: Float
)