package com.example.gdsudao.activity

import android.Manifest
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


        var estacao : Estacao

        var latitude = ptLatitude.text.toString()
        var longitude = ptLongitude.text.toString()

        Toast.makeText(this@CadastroAreaActivity, "" + latitude + " - " + longitude, Toast.LENGTH_LONG).show()

        btnBuscarPorCoordenadas.setOnClickListener {
                // Buscando por coordenadas
                val requestCall = RetrofitInitializer().apiService().getEstacao(latitude, longitude)

                var localizacao = listOf<Float>(latitude.toFloat(), longitude.toFloat())

                requestCall.enqueue(object : Callback<Estacao> {

                    override fun onResponse(call: Call<Estacao>, response: Response<Estacao>) {
                        if (response.isSuccessful) {
                            estacao = Estacao(
                                response.body()._id.toString(),
                                response.body().altitude?.toFloat(),
                                response.body().codigoCPTEC.toString(),
                                response.body().codigoINMET.toString(),
                                Localizacao(localizacao),
                                response.body().nomeEstacao.toString(),
                                response.body().tipoEstacao.toString(),
                                response.body().uf.toString()
                            )
                            var nomeDaLocalizacao = estacao.nomeEstacao
                            ptLocalizacao.setText(nomeDaLocalizacao)

                            // TODO: salvar estacao
                            var sharedPreferences = getSharedPreferences("gdsudao", MODE_PRIVATE)
                            //var gson = Gson()
                            //var json = gson.toJson(estacao);
                            sharedPreferences.edit().putString("nome", "Ricardo")
                            sharedPreferences.edit().commit()

                            /* TODO Testar ^^, recuperando dados
                            var sharedPreferences = getSharedPreferences("gdsudao", MODE_PRIVATE)
                            var gson = Gson()
                            var jsonG : String? = sharedPreferences.getString("estacao", " ")
                            Toast.makeText(this@CadastroAreaActivity, "" + jsonG, Toast.LENGTH_SHORT).show()
                            var estacao : Estacao = gson.fromJson(jsonG, Estacao::class.java)
                            Toast.makeText(this@CadastroAreaActivity, "" + estacao, Toast.LENGTH_SHORT).show()
                            */

                            Toast.makeText(this@CadastroAreaActivity, "" + estacao, Toast.LENGTH_LONG).show()
                        }else{
                            Toast.makeText(this@CadastroAreaActivity, "ERRO", Toast.LENGTH_SHORT).show()
                        }
                    }

                    override fun onFailure(call: Call<Estacao>?, t: Throwable?) {
                        Toast.makeText(this@CadastroAreaActivity, "ERRO2:" + t.toString() , Toast.LENGTH_SHORT).show()
                    }
                })


        }



        btnUsarLocalizacaoAtual.setOnClickListener {

        }

    }
}