package com.example.gdsudao.activity

import android.content.Intent
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.view.View
import android.widget.ListAdapter
import android.widget.Toast
import androidx.recyclerview.widget.LinearLayoutManager
import com.example.gdsudao.R
import com.example.gdsudao.adapter.AreaAdapter
import com.example.gdsudao.model.*
import com.example.gdsudao.utils.RetrofitInitializer
import com.google.gson.Gson
import kotlinx.android.synthetic.main.activity_cadastro_area.*
import kotlinx.android.synthetic.main.activity_menu.*
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

class MenuActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_menu)

        setSupportActionBar(findViewById(R.id.toolbarMenu))

        //getStation("-54.013292", "-31.347801") ok
        //getNormais("BAGE") ok
        //getDiarios("A827","2020-09-10", "2020-09-15") ok
        //getPrevisoes("A827","2020-09-30") ok

        //** Recycler view **//

        // Configurar adapter
        /*
        val areas = listOf(
            Area("p20", "10/10/2000", "10/10/2000", 1, 100.0f),
            Area("p20 - Talhão 2 ", "10/10/2000", "10/10/2000", 2, 50.0f),
            Area("Aceguá 1", "10/10/2000", "10/10/2000", 4, 69.0f),
            Area("Embrapa", "10/10/2000", "10/10/2000", 2, 10.5f),
            Area("p20", "10/10/2000", "10/10/2000", 1, 100.0f),
            Area("p20 - Talhão 2 ", "10/10/2000", "10/10/2000", 2, 50.0f),
            Area("Aceguá 1", "10/10/2000", "10/10/2000", 4, 69.0f),
            Area("Embrapa", "10/10/2000", "10/10/2000", 2, 10.5f),
            Area("p20", "10/10/2000", "10/10/2000", 1, 100.0f),
            Area("p20 - Talhão 2 ", "10/10/2000", "10/10/2000", 2, 50.0f),
            Area("Aceguá 1", "10/10/2000", "10/10/2000", 4, 69.0f),
            Area("Embrapa", "10/10/2000", "10/10/2000", 2, 10.5f)
        )

        */

        //var areas = listOf(Area("teste", "teste"))
        var sp = com.example.gdsudao.utils.SharedPreferences()
        var areas = sp.RecuperarListaAreas(this)
        recyclerViewAreas.apply {
            layoutManager = LinearLayoutManager(this@MenuActivity)
            adapter = AreaAdapter(areas)
            hasFixedSize()
        }
        /* Configurar o Recyclerviwer
        var viewManager = LinearLayoutManager(this)
        recyclerViewAreas.layoutManager = viewManager
        recyclerViewAreas.hasFixedSize()
        recyclerViewAreas.adapter = viewAdapter
         */

        // Vai para o cadastro de Área
        btnAddArea.setOnClickListener{
            var intent = Intent(this, CadastroAreaActivity::class.java)
            startActivity(intent)
        }









    }





    fun getNormais(nomeEstacao: String){

        val requestCall = RetrofitInitializer().apiService().getNormais(nomeEstacao)

        requestCall.enqueue(object : retrofit2.Callback<Normal> {

            override fun onResponse(call: Call<Normal>, response: Response<Normal>) {
                if (response.isSuccessful) {
                    var normais = response.body()
                    Toast.makeText(this@MenuActivity, "DENTRO" + normais, Toast.LENGTH_LONG).show()
                }else{
                    Toast.makeText(this@MenuActivity, "ERRO1", Toast.LENGTH_SHORT).show()
                }
            }

            override fun onFailure(call: Call<Normal>?, t: Throwable?) {
                Toast.makeText(this@MenuActivity, "ERRO2:" + t.toString() , Toast.LENGTH_SHORT).show()
            }
        })

    }

    fun getDiarios(codigoINMET: String, dataInicial: String, dataFinal: String){

        val requestCall = RetrofitInitializer().apiService().getDiarios(codigoINMET,dataInicial, dataFinal)

        requestCall.enqueue(object : retrofit2.Callback<List<Diario>> {

            override fun onResponse(call: Call<List<Diario>>, response: Response<List<Diario>>) {
                if (response.isSuccessful) {
                    var diarios = response.body()
                    Toast.makeText(this@MenuActivity, "DENTRO" + diarios, Toast.LENGTH_LONG).show()
                }else{
                    Toast.makeText(this@MenuActivity, "ERRO1", Toast.LENGTH_SHORT).show()
                }
            }

            override fun onFailure(call: Call<List<Diario>>?, t: Throwable?) {
                Toast.makeText(this@MenuActivity, "ERRO2:" + t.toString() , Toast.LENGTH_SHORT).show()
            }
        })

    }

    fun getPrevisoes(codigoINMET: String, dataInicial: String){

        val requestCall = RetrofitInitializer().apiService().getPrevisoes(codigoINMET,dataInicial)

        requestCall.enqueue(object : retrofit2.Callback<List<Previsao>> {

            override fun onResponse(call: Call<List<Previsao>>, response: Response<List<Previsao>>) {
                if (response.isSuccessful) {
                    var previsoes = response.body()
                    Toast.makeText(this@MenuActivity, "DENTRO" + previsoes, Toast.LENGTH_LONG).show()
                }else{
                    Toast.makeText(this@MenuActivity, "ERRO1", Toast.LENGTH_SHORT).show()
                }
            }

            override fun onFailure(call: Call<List<Previsao>>?, t: Throwable?) {
                Toast.makeText(this@MenuActivity, "ERRO2:" + t.toString() , Toast.LENGTH_SHORT).show()
            }
        })

    }

}