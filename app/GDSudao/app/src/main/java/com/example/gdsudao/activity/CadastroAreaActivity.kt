package com.example.gdsudao.activity

import android.Manifest
import android.content.Intent
import android.content.SharedPreferences
import android.content.pm.PackageManager
import android.location.Location
import android.location.LocationManager
import android.os.Bundle
import android.widget.Toast
import com.google.android.material.floatingactionbutton.FloatingActionButton
import com.google.android.material.snackbar.Snackbar
import androidx.appcompat.app.AppCompatActivity
import androidx.core.app.ActivityCompat
import com.example.gdsudao.R
import com.example.gdsudao.model.Area
import com.example.gdsudao.model.Estacao
import com.example.gdsudao.model.Localizacao
import com.example.gdsudao.utils.RetrofitInitializer
import com.google.android.gms.location.LocationServices
import com.google.gson.Gson
import com.google.gson.internal.`$Gson$Preconditions`
import com.google.gson.reflect.TypeToken
import kotlinx.android.synthetic.main.activity_cadastro_area.*
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

class CadastroAreaActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_cadastro_area)

        setSupportActionBar(findViewById(R.id.toolbarCadastroArea))

        // so para remover a cagada
        //var sp = com.example.gdsudao.utils.SharedPreferences()
        //sp.RemoverAllAreaLista(this)


        btnAddLocalizacao.setOnClickListener {
            var sp = com.example.gdsudao.utils.SharedPreferences()
            sp.RemoverAllAreaLista(this)
        }


        btnCadastrarAreas.setOnClickListener {

            var sp = com.example.gdsudao.utils.SharedPreferences()
            var nomeArea = etNomeArea.text.toString()
            var dataCorte = etDataUltimoCorte.text.toString()
            var numeroCortes = etNumeroCortes.text.toString()
            var latitude = etLatitude.text.toString()
            var longitude = etLongitude.text.toString()


            //intent.putExtra("nomeArea", nomeArea)
            //intent.putExtra("dataCorte", dataCorte)
            //intent.putExtra("numeroCortes", numeroCortes)
            // intent.putExtra("latitude", latitude)
            // intent.putExtra("longitude", longitude)

            // Salvar na fila de areas

            var intent = Intent(this, MenuActivity::class.java)
            var area = Area(nomeArea, dataCorte, numeroCortes, latitude, longitude, "", "", "", "", "", "")

            // Consultar estacao
            val requestCall = RetrofitInitializer().apiService().getEstacao(latitude, longitude)
            requestCall.enqueue(object : Callback<Estacao> {

                override fun onResponse(call: Call<Estacao>, response: Response<Estacao>) {
                    if (response.isSuccessful) {
                        var codigoINMET = response.body().codigoINMET.toString()
                        area.codigoEstacao = codigoINMET
                        Toast.makeText(applicationContext, codigoINMET, Toast.LENGTH_LONG).show()
                        sp.SalvarAreaLista(applicationContext, area)

                        startActivity(intent)
                    }else{
                        //TODO
                    }
                }

                override fun onFailure(call: Call<Estacao>?, t: Throwable?) {
                    Toast.makeText(this@CadastroAreaActivity, "ERRO2:" + t.toString() , Toast.LENGTH_SHORT).show()
                }
            })


        }
    }
}