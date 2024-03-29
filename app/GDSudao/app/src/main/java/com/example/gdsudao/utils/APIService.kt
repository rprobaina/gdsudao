package com.example.gdsudao.utils

import com.example.gdsudao.model.*
import retrofit2.http.GET
import retrofit2.http.Path

interface APIService {
    // Retorna a estação mais próxima
    @GET("/estacao/maisproxima/{latitude}/{longitude}")
    fun getEstacao(@Path("latitude") latitude: String, @Path("longitude") longitude: String): retrofit2.Call<Estacao>

    // Retorna as normais
    @GET("/normais/{nomeEstacao}")
    fun getNormais(@Path("nomeEstacao") nomeEstacao: String): retrofit2.Call<Normal>

    // Retorna os dados diários
    @GET("/diarios/{codigoINMET}/{dataInicial}/{dataFinal}")
    fun getDiarios(@Path("codigoINMET") codigoINMET: String, @Path("dataInicial") dataInicial: String,
                   @Path("dataFinal") dataFinal: String): retrofit2.Call<List<Diario>>

    // Retorna os dados de previsao do tempo
    @GET("/previsoes/{codigoINMET}/{dataInicial}")
    fun getPrevisoes(@Path("codigoINMET") codigoINMET: String, @Path("dataInicial") dataInicial: String): retrofit2.Call<List<Previsao>>

    // Retorna os dados de previsao do tempo
    @GET("/proximoPastejo/{codigoINMET}/{dataCorte}/{numeroCorte}")
    fun attArea(@Path("codigoINMET") codigoEstcao: String, @Path("dataCorte") dataCorte: String, @Path("numeroCorte") numeroCorte: String):  retrofit2.Call<Area>
}