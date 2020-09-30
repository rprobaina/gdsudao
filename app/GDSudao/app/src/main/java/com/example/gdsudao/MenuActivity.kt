package com.example.gdsudao

import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import android.widget.Toast
import com.example.gdsudao.model.Estacao
import retrofit2.Call
import retrofit2.Response

class MenuActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_menu)

        getStation()
    }

    fun getStation(){

        val requestCall = RetrofitInitializer().estacaoService().getEstacao("-54.013292", "-31.347801")

        requestCall.enqueue(object : retrofit2.Callback<Estacao> {
            override fun onResponse(call: Call<Estacao>, response: Response<Estacao>) {
                Toast.makeText(this@MenuActivity, "" + response.body(), Toast.LENGTH_SHORT).show()
                if (response.isSuccessful) {
                    var estacao: Estacao? = response.body()
                    Toast.makeText(this@MenuActivity, "" + estacao?.codigoCPTEC, Toast.LENGTH_SHORT).show()
                }else{
                    Toast.makeText(this@MenuActivity, "Erro " , Toast.LENGTH_LONG).show()
                }
            }

            override fun onFailure(call: Call<Estacao>?, t: Throwable?) {
                Log.println(Log.ERROR, "estacao", call.toString())
                Log.println(Log.ERROR, "estacao", t.toString())
                Toast.makeText(this@MenuActivity, "" + t.toString() , Toast.LENGTH_SHORT).show()
            }
        })

       //Log.println(Log.ERROR, "estacao", estacao.nomeEstacao)
    }


}