package com.example.gdsudao.utils

import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory

class RetrofitInitializer {


    private val retrofit = Retrofit.Builder()
            .baseUrl("http://192.168.0.106:8082")
            .addConverterFactory(GsonConverterFactory.create())
            .build()

    fun apiService(): APIService {
        return retrofit.create(APIService::class.java)
    }
}