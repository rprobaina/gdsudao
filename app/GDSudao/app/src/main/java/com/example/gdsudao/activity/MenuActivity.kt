package com.example.gdsudao.activity

import android.content.Intent
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
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


        //var areas = listOf(Area("teste", "teste"))
        var sp = com.example.gdsudao.utils.SharedPreferences()
        //sp.RemoverAllAreaLista(this) //Se eu cadastrar uma are errada liberar isso
        var areas = sp.RecuperarListaAreas(this)
        Log.println(Log.DEBUG, "areas", areas.toString())

        if (areas.size > 0) {
            areas.forEach(){
                //Toast.makeText(this, it.toString(), Toast.LENGTH_LONG).show()
                // Consultar estacao
                val requestCall = RetrofitInitializer().apiService().attArea(it.codigoEstacao, it.dataCorte, it.numeroCorte)
                //Toast.makeText(this, requestCall.toString(), Toast.LENGTH_SHORT).show()
                requestCall.enqueue(object : Callback<Area> {
                    override fun onResponse(call: Call<Area>, response: Response<Area>) {
                        if (response.isSuccessful) {
                            var areaResponse = response.body()
                            Toast.makeText(this@MenuActivity, areaResponse.toString(), Toast.LENGTH_SHORT).show()
                            it.st = areaResponse.st
                            it.proxcorte = areaResponse.proxcorte
                            it.diario = areaResponse.diario
                            it.previsao = areaResponse.previsao
                            it.normal = areaResponse.normal
                            Log.println(Log.DEBUG, "areas", response.toString())

                            // TODO: salvar area atualizada localmente

                            recyclerViewAreas.apply {
                                layoutManager = LinearLayoutManager(this@MenuActivity)
                                adapter = AreaAdapter(this@MenuActivity, areas)
                                hasFixedSize()
                            }

                        }else{
                            //TODO
                        }
                    }

                    // TODO tratar o erro se não não cria a lista
                    override fun onFailure(call: Call<Area>?, t: Throwable?) {
                        Log.println(Log.ERROR, "onFailure", "falhou a resposta")
                    }
                })
            }
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



   fun attAreas(areas: ArrayList<Area>) {

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