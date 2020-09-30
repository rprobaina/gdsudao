package com.example.gdsudao

import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory

class RetrofitInitializer {


    private val retrofit = Retrofit.Builder()
            .baseUrl("http://192.168.0.101:8082")
            .addConverterFactory(GsonConverterFactory.create())
            .build()



    fun estacaoService(): EstacaoService {
        return retrofit.create(EstacaoService::class.java)
    }
}