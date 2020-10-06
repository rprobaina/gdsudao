package com.example.gdsudao.model

data class Previsao (
    var _id: String,
    var clima: String,
    var codCPTEC: String,
    var codINMET: String,
    var dataAtualizacao: String,
    var dataPrevisao: String,
    var iuv: Int,
    var temperaturaMaxima: Float,
    var temperaturaMinima: Float
)