package com.example.gdsudao.activity

import android.Manifest
import android.app.DatePickerDialog
import android.content.Intent
import android.content.pm.PackageManager
import android.location.Location
import android.os.Bundle
import android.text.Editable
import android.util.Log
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.core.app.ActivityCompat
import com.example.gdsudao.utils.FusedLocationProvider
import com.example.gdsudao.R
import com.example.gdsudao.model.Area
import com.example.gdsudao.model.Estacao
import com.example.gdsudao.utils.RetrofitInitializer
import com.google.android.gms.location.LocationServices
import kotlinx.android.synthetic.main.activity_cadastro_area.*
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response
import java.text.SimpleDateFormat
import java.util.*

class CadastroAreaActivity : AppCompatActivity() {

    var nomeArea = ""
    var dataCorte = ""
    var numeroCortes = ""
    var latitude = ""
    var longitude = ""
    var dataFmt = ""

    override fun onResume() {
        super.onResume()

        val bundle = intent.extras

        if (bundle != null){
            nomeArea = bundle?.getString("nA").toString()
            dataCorte = bundle?.getString("dC").toString()
            numeroCortes = bundle?.getString("nC").toString()
            latitude = bundle?.getString("lat").toString()
            longitude = bundle?.getString("lon").toString()
        }

        etNomeArea.text = Editable.Factory.getInstance().newEditable(nomeArea)
        etDataUltimoCorte.text = Editable.Factory.getInstance().newEditable(dataCorte)
        etNumeroCortes.text = Editable.Factory.getInstance().newEditable(numeroCortes)
        etLatitude.text = Editable.Factory.getInstance().newEditable(latitude)
        etLongitude.text = Editable.Factory.getInstance().newEditable(longitude)
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_cadastro_area)

        setSupportActionBar(findViewById(R.id.toolbarCadastroArea))

        var fusedLocationClient = LocationServices.getFusedLocationProviderClient(this)

        var dataFmt = ""


        val bundle = intent.extras

        if (bundle != null){
            nomeArea = bundle?.getStringArrayList("nA").toString()
            dataCorte = bundle?.getStringArrayList("dC").toString()
            numeroCortes = bundle?.getStringArrayList("nC").toString()
            latitude = bundle?.getStringArrayList("lat").toString()
            longitude = bundle?.getStringArrayList("lon").toString()
        }

        etNomeArea.text = Editable.Factory.getInstance().newEditable(nomeArea)
        etDataUltimoCorte.text = Editable.Factory.getInstance().newEditable(dataCorte)
        etNumeroCortes.text = Editable.Factory.getInstance().newEditable(numeroCortes)
        etLatitude.text = Editable.Factory.getInstance().newEditable(latitude)
        etLongitude.text = Editable.Factory.getInstance().newEditable(longitude)


        var cal = Calendar.getInstance()

        val dateSetListener = DatePickerDialog.OnDateSetListener { view, year, monthOfYear, dayOfMonth ->
            cal.set(Calendar.YEAR, year)
            cal.set(Calendar.MONTH, monthOfYear)
            cal.set(Calendar.DAY_OF_MONTH, dayOfMonth)

            dataFmt = SimpleDateFormat("yyyy-MM-dd", Locale.US).format(cal.time)
            val sdf = SimpleDateFormat("dd/MM/yyyy", Locale.US)
            etDataUltimoCorte.text = Editable.Factory.getInstance().newEditable(sdf.format(cal.time))
        }

        etDataUltimoCorte.setOnClickListener {
            DatePickerDialog(this, dateSetListener,
                    cal.get(Calendar.YEAR),
                    cal.get(Calendar.MONTH),
                    cal.get(Calendar.DAY_OF_MONTH)).show()
        }


        btnCadastrarAreas.setOnClickListener {

            var sp = com.example.gdsudao.utils.SharedPreferences()

            var nomeArea = etNomeArea.text.toString()


            var numeroCortes = etNumeroCortes.text.toString()
            var latitude = etLatitude.text.toString()
            var longitude = etLongitude.text.toString()
            var dataCorte = etDataUltimoCorte.text.toString()

            if(validarDados(nomeArea, dataCorte, numeroCortes, latitude, longitude)){

                var intent = Intent(this, MenuActivity::class.java)
                if (dataFmt.isNullOrEmpty()){
                    dataFmt = "${dataCorte.substring(6,10)}-${dataCorte.substring(3,5)}-${dataCorte.substring(0,2)}"
                }
                println("Nome area: " + nomeArea)
                println("Data Corte: " + dataFmt)
                println("N Cortes: " + numeroCortes)
                println("Lat: " + latitude)
                println("Long: " + longitude)
                //println("Data Corte" + dataCorteFmt)
                var area = Area(nomeArea, dataFmt, numeroCortes, latitude, longitude, "", "", "", "", "", "")



                val requestCall = RetrofitInitializer().apiService().getEstacao(latitude, longitude)
                requestCall.enqueue(object : Callback<Estacao> {

                    override fun onResponse(call: Call<Estacao>, response: Response<Estacao>) {
                        if (response.isSuccessful) {
                            var codigoINMET = response.body().codigoINMET.toString()
                            area.codigoEstacao = codigoINMET

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
            }else{
                Toast.makeText(this, "Dados inválidos", Toast.LENGTH_SHORT).show()
            }



        }


        // TODO: so pega a localizacao se o google maps estiver aberto
        //  (ver: https://gist.github.com/Jthomas54/c1bdd68653a1832024c7cd4f784baf02)
        btnAddLocalizacaoGPS.setOnClickListener {


            if (ActivityCompat.checkSelfPermission(
                    this,
                    Manifest.permission.ACCESS_FINE_LOCATION
                ) != PackageManager.PERMISSION_GRANTED && ActivityCompat.checkSelfPermission(
                    this,
                    Manifest.permission.ACCESS_COARSE_LOCATION
                ) != PackageManager.PERMISSION_GRANTED
            ) {
                var missingPermissions = arrayOf(Manifest.permission.ACCESS_FINE_LOCATION, Manifest.permission.ACCESS_COARSE_LOCATION)
                ActivityCompat.requestPermissions(this, missingPermissions, 1000)
                return@setOnClickListener
            }



            fusedLocationClient.lastLocation.addOnSuccessListener { location: Location? ->
                    if (location != null) {
                        var lat = location.latitude.toString()
                        var lon = location.longitude.toString()

                        etLatitude.text = Editable.Factory.getInstance().newEditable(lat)
                        etLongitude.text = Editable.Factory.getInstance().newEditable(lon)
                    } else {
                        val flp = FusedLocationProvider(this)
                        flp.startUpdates()
                        val location = flp.getCurrentLocation()
                        if( location != null){
                            var lat = location.latitude.toString()
                            var lon = location.longitude.toString()

                            etLatitude.text = Editable.Factory.getInstance().newEditable(lat)
                            etLongitude.text = Editable.Factory.getInstance().newEditable(lon)
                        }else{
                            Toast.makeText(this, "Loacilização indisponível! Verifique se GPS está habilitado.", Toast.LENGTH_SHORT).show()
                        }
                    }
                }


        }


        btnAddLocalizacaoMAPA.setOnClickListener {
            //salva dados antes de abrir o mapa
            nomeArea = etNomeArea.text.toString()
            dataCorte = etDataUltimoCorte.text.toString()
            numeroCortes = etNumeroCortes.text.toString()
            latitude = etLatitude.text.toString()
            longitude = etLongitude.text.toString()

            var intent = Intent(this, MapsActivity::class.java)

            intent.putExtra("nA", nomeArea)
            intent.putExtra("dC", dataCorte)
            intent.putExtra("nC", numeroCortes)
            intent.putExtra("lat", latitude)
            intent.putExtra("lon", longitude)

            startActivity(intent)
        }

    }




    fun validarDados(a: String, b: String, c: String, d: String, e: String): Boolean{
        var r = true
        if (a.isNullOrEmpty() || b.isNullOrEmpty() || c.isNullOrEmpty() || d.isNullOrEmpty() || e.isNullOrEmpty()){
            r = false
        }
        return r
    }
}