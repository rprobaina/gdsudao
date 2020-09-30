package com.example.gdsudao

import com.example.gdsudao.model.Estacao
import retrofit2.http.GET
import retrofit2.http.Path


interface EstacaoService {

    @GET("/estacao/maisproxima/{latitude}/{longitude}")
    fun getEstacao(@Path("latitude") latitude: String, @Path("longitude") longitude: String): retrofit2.Call<Estacao>

}